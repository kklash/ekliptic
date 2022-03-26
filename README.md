# Ekliptic

This package provides primitives for elliptic curve cryptographic operations on the secp256k1 curve, with zero dependencies and excellent performance. It provides both Affine and Jacobian interfaces for elliptic curve operations. Aims to facilitate performant low-level operations on secp256k1 without overengineering or kitchen-sink syndrome.

## ALPHA STATE

This library is not finished, stable, or audited - depend on it at your own peril!

## Elliptic-whah?

Elliptic curve cryptography is a relatively new field of [asymmetric public-key cryptography](https://cryptobook.nakov.com/asymmetric-key-ciphers). An elliptic curve is just a cubic equation of a particular form. The secp256k1 curve, for example, is `y² = x³ + 7`. To make this curve equation useful, we first define an addition operation that 'adds' two `(x, y)` points on the curve to produce a third point _also_ on the curve. From that, you can create a multiplication operation to multiply a 2D point by some 1D (scalar) number, by simply adding the point to itself many times.

It just so happens that due to the particular properties of elliptic curves, if you multiply some publicly known point by a secret number, that operation is extremely hard to reverse, and you end up with another point that is mathematically related to the secret number. Functions that are easy to compute but hard to reverse are a fundamental building block of cryptography, and people started to realize you could use this feature of elliptic curve equations as a basis for new public-key cryptosystems, like RSA, but using much smaller numbers in a 2D space.

The unique one-way function of elliptic curve cryptography is _base point multiplication over a finite field,_ (the 'finite field' part means all coordinate values are taken modulo some large prime number). A base point is a publicly known `(x, y)` point, often called the _generator_ point `G`, which all parties agree upon. The private key in this cryptosystem is a scalar number `k` which is multiplied with the base point. The point resulting from base point multiplication `P` becomes the public key. `P` and `G` are capitalized to denote that they are 2D points, while `k` is a lone positive integer. Base point multiplication is written mathematically as

```
P = k * G
```

Point multiplication is _believed_ to be hard to undo: There's no way to quickly compute `k` if you only know `P` and `G`. The only known way to efficiently perform `P / G` would be to run Shor's Algorithm on a quantum computer that can operate with at least `6 * log2(k)` qubits. This currently doesn't exist, so at least for now, elliptic curve cryptography provides a safe way to sign/verify and encrypt/decrypt information asymetrically.


## Down Sides

Elliptic curve cryptography does have some down sides - Primarily, from the complexity involved in implementing it safely. To perform elliptic curve cryptography, _someone_ needs to design an elliptic curve with its various parameters in a secure way, which requires highly adept and experienced cryptographers. This makes users vulnerable to malicious design by those with the specialized knowledge needed to produce such curves. Compared to a simpler system like RSA, where there are no 'magic numbers' involved, ECC predicates the safety of the system not only on the security of the algorithms and in-code implementations, but also on the ethical integrity of curve designers, who are far fewer, and more tightly centralized.

Thankfully, the secp256k1 curve was designed in a non-random 'nothing up my sleeve' fashion, which helps to reduce the risk that it was designed with a backdoor in mind. This is why it has become such a popular curve. Satoshi chose to use secp256k1 for Bitcoin for that same reason (among others).

## Why not RSA?

Primarily, for performance. Elliptic curves offer a way to perform cryptography faster for the same degree of security.

A 256-bit elliptic curve key provides roughly the same degree of security as a 2048-bit RSA key. But for normal 'happy path' operations where you're not trying attack the cryptosystem, elliptic curve operations are _vastly_ faster, simply due to the size of the numbers involved. It's easier to multiply `5 x 9` than to multiply `555 x 999`.

Consider this simple benchmark which compares 256-bit ECC and 2048-bit RSA private and public key generation:

```go
func BenchmarkGenerateKeys_Ekliptic(b *testing.B) {
  var privateKey *big.Int
  var publicKey struct{ x, y big.Int }
  for i := 0; i < b.N; i++ {
    privateKey, _ = ekliptic.NewPrivateKey(rand.Reader)
    ekliptic.MultiplyBasePoint(privateKey, &publicKey.x, &publicKey.y)
  }
}

func BenchmarkGenerateKeys_RSA(b *testing.B) {
  for i := 0; i < b.N; i++ {
    rsa.GenerateKey(rand.Reader, 2048)
  }
}
```

```
goos: linux
goarch: amd64
pkg: t
cpu: Intel(R) Core(TM) i7-6500U CPU @ 2.50GHz
BenchmarkGenerateKeys_Ekliptic-4        1530      711021 ns/op    242393 B/op     2833 allocs/op
BenchmarkGenerateKeys_RSA-4                8   138256984 ns/op   1874108 B/op     5019 allocs/op
```

Golang's standard library `crypto/rsa` package takes between between 0.1 to 0.25 **seconds** to generate an RSA key pair of that size. Ekliptic takes less than a millisecond to generate an ECC key pair of equivalent security. For activities like web-browsing, ephemeral key pairs are constantly being generated to negotiate TLS-encrypted connections. The faster you can generate session keys and verify certificates, the faster you can open TLS connections, and the faster you can browse.

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

You can even safely pass the input pointers as the output pointers, to modify them in place.

```go
ekliptic.DoubleJacobi(
  x3, y3, z3,
  x3, y3, z3,
)
```

`big.Int` structs can be re-used when the values they hold are no longer required. This is why you'll see patterns like this if you read Ekliptic's code:

```go
e := a.Mul(a, three)
a = nil
```

In the above example, `a` is no longer needed, so we reclaim its memory as a new variable to avoid allocating an entirely new `big.Int` struct for `e`.

### Jacobian Points

This library offers support for both affine and Jacobian point math. Affine coordinates are 'normal' two-dimensional coordinates, `x` and `y`, which unambiguously describes a point on the plane. Jacobian coordinates are a three-dimensional representation of an affine point, (Ax, Ay), in terms of three variables: (x, y, z) such that:

```
Ax = x / z²
Ay = y / z³
```

This relationship means there are an absurdly large number of possible Jacobian coordinate triplets which describe the same affine point. Each affine coordinate is basically converted into a ratio of `x:z` and `y:z`, thus proportional ratios simplify to the same affine point.

Why would we want to represent points this way? Elliptic curve multiplication - a critical primitive for almost any elliptic-curve cryptography - involves performing many addition operations in a row. That's what multiplication means, after all. When you add two affine `(x, y)` points together in an elliptic curve, you have to perform some finite field division, AKA modular inversion, to get a result back in affine form. Modular inversion is a very expensive operation compared to multiplication. Instead of dividing after _every_ addition operation, you can defer the division until the end of the multiplication sequence, by accumulating in the divisor coordinate `z`. After the multiplication operation is done, the point can be converted back to affine, or used for new EC operations, as needed.

To demonstrate, notice how expensive a naive affine multiplication is compared to a Jacobian multiplication:

```
BenchmarkMultiplyJacobi-6                      675     1757329 ns/op    727070 B/op     5060 allocs/op
BenchmarkMultiplyAffine-6                      679     1782691 ns/op    728819 B/op     5084 allocs/op
BenchmarkMultiplyAffineNaive-6                 442     2480711 ns/op    545915 B/op     9147 allocs/op
```
`ekliptic.MultiplyJacobi` and `ekliptic.MultiplyAffine` both use Jacobian math for multiplication operations under the hood. `ekliptic.MultiplyAffineNaive` is a naive implementation which uses affine addition and doubling instead of Jacobian math. It should be used for demonstrative purposes only.


### Precomputation

You can improve multiplication performance even more by using precomputed doubles of the secp256k1 base-point. Precomputing `G * 2^i` for `i` in `[0..256]` significantly boosts performance for base-point [double-and-add multiplication](https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add), especially if [the precomputed doubles are saved in affine form](./precomputed_doubles.go). Values are computed using the `ekliptic.ComputePointDoubles` function, [triggered by `go generate`](./genprecompute).

### Other Performance Notes

- We have [a special implementation which checks for Jacobi point validity without costly affine conversion.](./is_on_curve.go)
- Golang's `big.Int` has some operations which are more costly than others. For example, doing `n.Exp(n, two, P)` is more costly than doing `n.Mul(n, n); mod(n)`. This holds for squaring and cubing, but exponents beyond 3 require `.Exp` for the best performance.

### Thanks!

I wrote this library as a learning exercise. Thanks a ton to the following people and their amazing guides, with which I followed along and learned from.

- [paulmillr](https://paulmillr.com/posts/noble-secp256k1-fast-ecc/)
- [Andrea Corbellini](https://andrea.corbellini.name/2015/05/17/elliptic-curve-cryptography-a-gentle-introduction/)
- [Fang-Pen Lin](https://fangpenlin.com/posts/2019/10/07/elliptic-curve-cryptography-explained/)
- [Svetlin Nakov](https://cryptobook.nakov.com/asymmetric-key-ciphers/elliptic-curve-cryptography-ecc)
- [Nick Sullivan @ Cloudflare](https://blog.cloudflare.com/a-relatively-easy-to-understand-primer-on-elliptic-curve-cryptography/)

Many [test vectors in this package](./test_vectors/) were either duplicated from, generated by, or derived from other Elliptic Curve Crypto implementations:

- [`paulmillr/noble-secp256k1`](https://github.com/paulmillr/noble-secp256k1)
- [`btcsuite/btcec`](https://github.com/btcsuite/btcd/tree/master/btcec) (note: they use a different algorithm to add jacobians, but it works out to the same affine coordinates at the end. I modified a few of their test fixtures' jacobian ratios.)


## Donations

If you're interested in supporting development of this package, show your love by dropping me some Bitcoins!

### `bc1qhct3hwt5pjmu75d2fldwd477vhwmthuqvmh03s`
