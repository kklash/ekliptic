## Ekliptic

This package provides primitives for cryptographic operations on the secp256k1 curve, with zero dependencies and excellent performance. It provides both Affine and Jacobian interfaces for elliptic curve operations. Aims to facilitate performant low-level operations on secp256k1 without overengineering or kitchen-sink syndrome.

## ALPHA STATE

This library is not finished, stable, or audited - depend on it at your own peril!

## Examples

Deriving a public key from a private key:

```go
privateKey, _ := new(big.Int).SetString("c370af8c091812ef7f6bfaffb494b1046fb25486c9873243b80826daef3ec583", 16)
x := new(big.Int)
y := new(big.Int)

ekliptic.MultiplyBasePoint(privateKey, x, y)

fmt.Println("Public key:")
fmt.Printf(" x: %x\n", x)
fmt.Printf(" y: %x\n", y)

// output:
// Public key:
//  x: 76cd66c6cca75278ff408ce67290537367719154ae2b96448327fe4033ddcfc7
//  y: 35663ecbb64397bb9bd79155a1e6b138c2fb8fa1f11355f8e9e97ddd88a78e49
```

Deriving a Diffie-Hellman shared-secret:

```go
alice, _ := new(big.Int).SetString("94a22a406a6977c1a323f23b9d7678ad08e822834d1df8adece84e30f0c25b6b", 16)
bob, _ := new(big.Int).SetString("55ba19100104cbd2842999826e99e478efe6883ac3f3a0c7571034321e0595cf", 16)

var alicePub, bobPub struct{ x, y big.Int }

// derive public keys
MultiplyBasePoint(alice, &alicePub.x, &alicePub.y)
MultiplyBasePoint(bob, &bobPub.x, &bobPub.y)

var yValueIsUnused big.Int

// Alice gives Bob her public key, Bob derives the secret
bobSharedKey := new(big.Int)
MultiplyAffine(&alicePub.x, &alicePub.y, bob, bobSharedKey, &yValueIsUnused, nil)

// Bob gives Alice his public key, Alice derives the secret
aliceSharedKey := new(big.Int)
MultiplyAffine(&bobPub.x, &bobPub.y, alice, aliceSharedKey, &yValueIsUnused, nil)

fmt.Printf("Alice's derived secret: %x\n", aliceSharedKey)
fmt.Printf("Bob's derived secret:   %x\n", bobSharedKey)

// output:
// Alice's derived secret: 375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
// Bob's derived secret:   375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
```

## Hacking on Ekliptic

| Command | Usage |
|---------|-------|
| `go test` | Run unit tests. |
| `go test -bench=.` | Run benchmarks. |
| `go test -bench=. -benchmem` | Run benchmarks with memory profile. |
| `go generate` | Regenerate [precomputed base point doubles](./precomputed_doubles.go). |

## Performance Optimizations

### Memory

All methods use and accept golang-native `big.Int` structs for math operations. In most cases we require the caller to pass pointers in which results will be stored.

```go
x := new(big.Int).SetString("3e61e957cb7eb9252155722d37056b581cacd9949cd7daeba682d81ee829826d", 16)
y := new(big.Int).SetString("0c9b31d4b3f13c2dcec21b5d446a06cd655056d83495f63f05135ff4434e7ba5", 16)
z := new(big.Int).SetString("4c4619154810c1c0daa4ddd8c73971d159db91705f2113ce51b9885e4578874d", 16)

doubleX := new(big.Int)
doubleY := new(big.Int)
doubleZ := new(big.Int)

ekliptic.DoubleJacobi(
  x3, y3, z3,
  doubleX, doubleY, doubleZ,
)
```

While slightly more awkward to use, exposing this C-style API allows for better memory performance when doing large numbers of sequential operations. The garbage collector isn't doing as much work, because we don't have to keep re-allocating new `big.Int`s every time a call returns. We have some benchmarks for addition and doubling which demonstrate this method can save about 6 allocs/op, and a few hundred bytes of memory for every call involving Jacobian points.

`big.Int` structs can be re-used when the values they hold are no longer required. This is why you'll see patterns like this if you read Ekliptic's code:

```go
e := a.Mul(a, three)
a = nil
mod(e)
```

In the above example, `a` is no longer needed, so we reclaim its memory as a new variable to avoid allocating an entirely new `big.Int` struct for `e`.

### Jacobian Points

This library offers support for both affine and Jacobian point math. Affine coordinates are 'normal' two-dimensional coordinates, `x` and `y`, which unambiguously describes a point on the plane. Jacobian coordinates are a three-dimensional representation of an affine point, (Ax, Ay), in terms of three variables: (x, y, z) such that:

```
Ax = x / z²
Ay = y / z³
```

This relationship means there are an absurdly large number of possible Jacobian coordinate triplets which describe the same affine point.

Why would we want to represent points this way? Elliptic curve multiplication - a critical primitive for almost any elliptic-curve cryptography - involves performing many addition operations in a row. That's what multiplication means, after all. When you add two affine `(x, y)` points together in an elliptic curve, you have to perform some finite field division, AKA modular inversion, to get a result back in affine form. Modular inversion is a very expensive operation compared to multiplication. Instead of dividing after _every_ addition operation, you can defer the division until the end of the multiplication sequence, by accumulating in the divisor coordinate `z`. After the multiplication operation is done, the point can be converted back to affine, or used for new EC operations, as needed.

### Precomputation

You can improve multiplication performance even more by using precomputed doubles of the secp256k1 base-point. Precomputing `G * 2^i` for `i` in `[0..256]` significantly boosts performance for base-point [double-and-add multiplication](https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add), especially if [the precomputed doubles are saved in affine form](./precomputed_doubles.go). Values are computed using the `ekliptic.ComputePointDoubles` function, [triggered by `go generate`](./genprecompute).

### Other Performance Notes

- We have [a special implementation which checks for Jacobi point validity without costly affine conversion.](./is_on_curve.go)
- Golang's `big.Int` has some operations which are more costly than others. For example, doing `n.Exp(n, two, P)` is more costly than doing `n.Mul(n, n); mod(n)`. This holds for squaring and cubing, but exponents beyond 3 require `.Exp` for the best performance.

## TODO

- [ ] in EC math unit tests, check affine equality against vectors rather than jacobian equality.
- [ ] get better more official test vectors
- [ ] fulfill elliptic.Curve interface
- [ ] add more examples
- [ ] Subtract() func?
- [ ] make separate package name for tests (ekliptic_test)

### Thanks!

I wrote this library as a learning exercise. Thanks a ton to the following people and their amazing guides, with which I followed along and learned from.

- [paulmillr](https://paulmillr.com/posts/noble-secp256k1-fast-ecc/)
- [Andrea Corbellini](https://andrea.corbellini.name/2015/05/17/elliptic-curve-cryptography-a-gentle-introduction/)
- [Fang-Pen Lin](https://fangpenlin.com/posts/2019/10/07/elliptic-curve-cryptography-explained/)
- [Svetlin Nakov](https://cryptobook.nakov.com/asymmetric-key-ciphers/elliptic-curve-cryptography-ecc)
- [Nick Sullivan @ Cloudflare](https://blog.cloudflare.com/a-relatively-easy-to-understand-primer-on-elliptic-curve-cryptography/)

Many [test vectors in this package](./test-vectors/) were either duplicated from, generated by, or derived from other Elliptic Curve Crypto implementations:

- [`paulmillr/noble-secp256k1`](https://github.com/paulmillr/noble-secp256k1)
- [`btcsuite/btcec`](https://github.com/btcsuite/btcd/tree/master/btcec) (note: they use a different algorithm to add jacobians, but it works out to the same affine coordinates at the end. I modified a few of their test fixtures' jacobian ratios.)
