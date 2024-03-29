# Ekliptic

This package provides primitives for elliptic curve cryptographic operations on the secp256k1 curve, with zero dependencies and excellent performance. It provides both Affine and Jacobian interfaces for elliptic curve operations. Aims to facilitate performant low-level operations on secp256k1 without overengineering or kitchen-sink syndrome.

## ALPHA STATE

This library is not finished, stable, or audited - depend on it at your own peril!

### [API Reference](https://pkg.go.dev/github.com/kklash/ekliptic)

## Elliptic-whah?

Elliptic curve cryptography is a relatively new field of [asymmetric public-key cryptography](https://cryptobook.nakov.com/asymmetric-key-ciphers). An elliptic curve is just a cubic equation of a particular form. The secp256k1 curve, for example, is $y^{2}=x^3+7$. To make this curve equation useful, we first define an addition operation that 'adds' two $x,y$ points on the curve to produce a third point _also_ on the curve. From that, you can create a multiplication operation to multiply a 2D point by some 1D (scalar) number, by simply adding the point to itself many times.

It just so happens that due to the particular properties of elliptic curves, if you multiply some publicly known point by a secret number, that operation is extremely hard to reverse, and you end up with another point that is mathematically related to the secret number. Functions that are easy to compute but hard to reverse are a fundamental building block of cryptography, and people started to realize you could use this feature of elliptic curve equations as a basis for new public-key cryptosystems, like RSA, but using much smaller numbers in a 2D space.

The unique one-way function of elliptic curve cryptography is _base point multiplication over a finite field,_ (the 'finite field' part means all coordinate values are taken modulo some large prime number). A base point is a publicly known $(x, y)$ point, often called the _generator_ point $G$, which all parties agree upon. The private key in this cryptosystem is a scalar number $k$ which is multiplied with the base point. The point resulting from base point multiplication $P$ becomes the public key. $P$ and $G$ are capitalized to denote that they are 2D points, while $k$ is a lone positive integer. Base point multiplication is written mathematically as

$$P = k G$$

Point multiplication is _believed_ to be hard to undo: There's no way to quickly compute $k$ if you only know $P$ and $G$. The only known way to efficiently perform $P \div G$ would be to run Shor's Algorithm on a quantum computer that can operate with at least $6\cdot \log_2(k)$ qubits. This currently doesn't exist, so at least for now, elliptic curve cryptography provides a safe way to sign/verify and encrypt/decrypt information asymetrically.


## Down Sides

Elliptic curve cryptography does have some down sides - Primarily, from the complexity involved in implementing it safely. To perform elliptic curve cryptography, _someone_ needs to design an elliptic curve with its various parameters in a secure way, which requires highly adept and experienced cryptographers. This makes users vulnerable to malicious design by those with the specialized knowledge needed to produce such curves. Compared to a simpler system like RSA, where there are no 'magic numbers' involved, ECC predicates the safety of the system not only on the security of the algorithms and in-code implementations, but also on the ethical integrity of curve designers, who are far fewer, and more tightly centralized.

Thankfully, the secp256k1 curve was designed in a non-random 'nothing up my sleeve' fashion, which helps to reduce the risk that it was designed with a backdoor in mind. This is why it has become such a popular curve. Satoshi chose to use secp256k1 for Bitcoin for that same reason (among others).

## Why not RSA?

Primarily, for performance. Elliptic curves offer a way to perform cryptography faster for the same degree of security.

A 256-bit elliptic curve key provides roughly the same degree of security as a 2048-bit RSA key. But for normal 'happy path' operations where you're not trying attack the cryptosystem, elliptic curve operations are _vastly_ faster, simply due to the size of the numbers involved. It's easier to multiply $5 \cdot 9$ than to multiply $555 \cdot 999$.

Consider this simple benchmark which compares 256-bit ECC and 2048-bit RSA private and public key generation:

```go
func BenchmarkGenerateKeys_Ekliptic(b *testing.B) {
  var privateKey *big.Int
  var publicKey struct{ x, y big.Int }
  for i := 0; i < b.N; i++ {
    privateKey, _ = ekliptic.RandomScalar(rand.Reader)
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

_Find fully validated examples in [`examples_test.go`](./examples_test.go)._

Deriving a public key from a private key:

```go
privateKey, _ := new(big.Int).SetString("c370af8c091812ef7f6bfaffb494b1046fb25486c9873243b80826daef3ec583", 16)
x, y := ekliptic.MultiplyBasePoint(privateKey)

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
alicePriv, _ := new(big.Int).SetString("94a22a406a6977c1a323f23b9d7678ad08e822834d1df8adece84e30f0c25b6b", 16)
bobPriv, _ := new(big.Int).SetString("55ba19100104cbd2842999826e99e478efe6883ac3f3a0c7571034321e0595cf", 16)

var alicePub, bobPub struct{ x, y *big.Int }

// derive public keys
alicePub.x, alicePub.y = ekliptic.MultiplyBasePoint(alicePriv)
bobPub.x, bobPub.y = ekliptic.MultiplyBasePoint(bobPriv)

// Alice gives Bob her public key, Bob derives the secret
bobSharedKey, _ := ekliptic.MultiplyAffine(alicePub.x, alicePub.y, bobPriv, nil)

// Bob gives Alice his public key, Alice derives the secret
aliceSharedKey, _ := ekliptic.MultiplyAffine(bobPub.x, bobPub.y, alicePriv, nil)

fmt.Printf("Alice's derived secret: %x\n", aliceSharedKey)
fmt.Printf("Bob's derived secret:   %x\n", bobSharedKey)

// output:
// Alice's derived secret: 375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
// Bob's derived secret:   375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
```

Signing and verifying a message with ECDSA.

```go
import (
  cryptorand "crypto/rand"
  mathrand "math/rand"

  "github.com/kklash/ekliptic"
)

randReader := mathrand.New(mathrand.NewSource(1))

key, _ := ekliptic.RandomScalar(randReader)

// This could also come from RFC6979 (github.com/kklash/rfc6979)
nonce, _ := cryptorand.Int(randReader, ekliptic.Secp256k1_CurveOrder)

hashedMessage := sha256.Sum256([]byte("i love you"))
hashedMessageInt := new(big.Int).SetBytes(hashedMessage[:])

r, s := ekliptic.SignECDSA(key, nonce, hashedMessageInt)

fmt.Printf("r: %x\n", r)
fmt.Printf("s: %x\n", s)

var pub struct{ x, y *big.Int }
pub.x, pub.y = ekliptic.MultiplyBasePoint(key)

valid := ekliptic.VerifyECDSA(hashedMessageInt, r, s, pub.x, pub.y)
fmt.Printf("valid: %v\n", valid)

// output:
//
// r: 4a821d5ec008712983929de448b8afb6c24e5a1b97367b9a65b6220d7f083fe3
// s: 381d053be61243d950865d7b8eb6b5ba48fbabfe7fda81af3183a184a02f5d51
// valid: true
```

Uncompressing a public key.

```go
compressedKey, _ := hex.DecodeString("030000000000000000000000000000000000000000000000000000000000000001")

publicKeyX := new(big.Int).SetBytes(compressedKey[1:])
evenY, oddY := ekliptic.Weierstrass(publicKeyX)

var publicKeyY *big.Int
if compressedKey[0]%2 == 0 {
  publicKeyY = evenY
} else {
  publicKeyY = oddY
}

fmt.Println("uncompressed key:")
fmt.Printf("x: %.64x\n", publicKeyX)
fmt.Printf("y: %.64x\n", publicKeyY)

// output:
// uncompressed key:
// x: 0000000000000000000000000000000000000000000000000000000000000001
// y: bde70df51939b94c9c24979fa7dd04ebd9b3572da7802290438af2a681895441
```

Ekliptic exports a struct type `Curve`, which satisfies the `elliptic.Curve` interface. You can use this in other libraries anywhere `elliptic.Curve` is accepted. For instance, to sign and verify data with `crypto/ecdsa`:

```go
d, _ := new(big.Int).SetString("18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725", 16)
pubX, pubY := ekliptic.MultiplyBasePoint(d)
key := &ecdsa.PrivateKey{
  D: d,
  PublicKey: ecdsa.PublicKey{
    Curve: new(ekliptic.Curve),
    X:     pubX,
    Y:     pubY,
  },
}

hashedMessage := sha256.Sum256([]byte("i love you"))

r, s, err := ecdsa.Sign(rand.Reader, key, hashedMessage[:])
if err != nil {
  panic("failed to compute signature: " + err.Error())
}

if ecdsa.Verify(&key.PublicKey, hashedMessage[:], r, s) {
  fmt.Println("verified ECDSA signature using crypto/ecdsa")
}

// output:
// verified ECDSA signature using crypto/ecdsa
```

Blinding a hidden value for multi-party computation:

```go
// Alice: input parameters
s, _ := new(big.Int).SetString("2fc9374cad648e33f78dd294578dd960281e05744b27faa1ffe1e7175bd6901d", 16)
aX, _ := new(big.Int).SetString("8bf20851cc16007dbf3df0c109dc016b360ca0f729f368ea38c385ceeffaf3cf", 16)
aY, _ := new(big.Int).SetString("a0c7bd73154c02cc5e002c3a4f876158a4276c185ef859df589675a92c745e3a", 16)

// Bob: input parameters
r, _ := new(big.Int).SetString("825e7984ae7843f9c13371d9a54143a465b1e2d278e67de1ca713127e40a52f1", 16)

// Alice: blind the input with secret s send the result B to Bob.
// B = s * A
bX, bY := ekliptic.MultiplyAffine(aX, aY, s, nil)

// Bob: receive A, blind it again with secret r and return result C to Alice.
// C = r * B = r * s * G
cX, cY := ekliptic.MultiplyAffine(bX, bY, r, nil)

// Alice: unblind C with the inverse of s:
// s⁻¹ * C = s⁻¹ * r * s * G = r * G
sInv := ekliptic.InvertScalar(s)

mX, mY := ekliptic.MultiplyAffine(cX, cY, sInv, nil)

expectedX, expectedY := ekliptic.MultiplyAffine(aX, aY, r, nil)

if !ekliptic.EqualAffine(mX, mY, expectedX, expectedY) {
  fmt.Println("did not find expected unblinded point")
  return
}

fmt.Println("found correct unblinded point q * h * G")

// output:
// found correct unblinded point q * h * G
```

## Hacking on Ekliptic

| Command | Usage |
|---------|-------|
| `go test` | Run unit tests. |
| `go test -bench=.` | Run benchmarks. |
| `go test -bench=. -benchmem` | Run benchmarks with memory profile. |
| `go generate` | Regenerate [precomputed base point products](./precomputed_table.go). |

## Test Vectors

Test vectors stored in [`test_vectors`](./test_vectors) can be verified using third-party libraries to double check that this library is working correctly.

| Validates against | Command to run |
|------------------|----------------|
|[`paritytech/libsecp256k1`](https://github.com/paritytech/libsecp256k1) (Rust)|`cargo run --manifest-path ./test_vectors/validate_rs/{Cargo.toml,}`|
|[`cslashm/ECPy`](https://github.com/cslashm/ECPy) (Python)|`pip3 install --user ECPy && python3 test_vectors/validate.py`|


## Performance Optimizations

### Memory

All methods use and accept golang-native `big.Int` structs for math operations. `big.Int` structs can be re-used when the values they hold are no longer required. This is why you'll see patterns like this if you read Ekliptic's code:

```go
e := a.Mul(a, three)
a = nil
```

In the above example, `a` is no longer needed, so we reclaim its memory as a new variable to avoid allocating an entirely new `big.Int` struct for `e`.

### Jacobian Points

This library offers support for both affine and Jacobian point math. Affine coordinates are 'normal' two-dimensional coordinates, $x$ and $y$, which unambiguously describes a point on the plane. Jacobian coordinates are a three-dimensional representation of an affine point, $x_a,y_a$, in terms of three variables: $x_j,y_j,z$ such that:

$$x_a = \frac{x_j}{z^2}$$

$$y_a = \frac{y_j}{z^3}$$

This relationship means there are an absurdly large number of possible Jacobian coordinate triplets which describe the same affine point. Each affine coordinate is basically converted into a ratio of $x:z$ and $y:z$, thus proportional ratios simplify to the same affine point.

Why would we want to represent points this way? Elliptic curve multiplication - a critical primitive for almost any elliptic-curve cryptography - involves performing many addition operations in a row. That's what multiplication means, after all. When you add two affine $(x, y)$ points together in an elliptic curve, you have to perform some finite field division, AKA modular inversion, to get a result back in affine form. Modular inversion is a very expensive operation compared to multiplication. Instead of dividing after _every_ addition operation, you can defer the division until the end of the multiplication sequence, by accumulating in the divisor coordinate $z$. After the multiplication operation is done, the point can be converted back to affine, or used for new EC operations, as needed.

To demonstrate, notice how expensive a naive affine multiplication is compared to a Jacobian multiplication:

```
BenchmarkMultiplyJacobi-6                      675     1757329 ns/op    727070 B/op     5060 allocs/op
BenchmarkMultiplyAffine-6                      679     1782691 ns/op    728819 B/op     5084 allocs/op
BenchmarkMultiplyAffineNaive-6                 442     2480711 ns/op    545915 B/op     9147 allocs/op
```
`ekliptic.MultiplyJacobi` and `ekliptic.MultiplyAffine` both use Jacobian math for multiplication operations under the hood. `ekliptic.MultiplyAffineNaive` is a naive implementation which uses affine addition and doubling instead of Jacobian math. It should be used for demonstrative purposes only.


### Precomputation

You can improve point multiplication performance significantly by precomputing a table of products of a point $P$ which you plan to multiply frequently. Precomputation means calculating $2^{4i}jP$ for $0 <= i <= 63$ and $0 <= j <= 15$. This table is indexable by $i$ and $j$. When using a precomputed table with a multiplication, `ekliptic` uses the [windowed method](https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Windowed_method). This speeds up multiplication of a fixed point often by a factor of 5x-10x, and saves a lot of memory allocations.

```
BenchmarkMultiplyJacobi-6                      241     4859994 ns/op   1229138 B/op     9382 allocs/op
BenchmarkMultiplyJacobi_Precomputed-6         2578      590604 ns/op    144339 B/op     1372 allocs/op
```

The secp256k1 base point $G$ is multiplied very frequently, so `ekliptic` [ships with a precomputed table for $G$ ready to go](./precomputed_table.go). This table is used to speed up `ekliptic.MultiplyBasePoint`. The source code file containing the table is auto-generated, [triggered by `go generate`](./genprecompute).

Precomputed tables for custom points can be constructed using `ekliptic.NewPrecomputedTable` function.

### Other Performance Notes

- We have [a special implementation which checks for Jacobi point validity without costly affine conversion.](./is_on_curve.go)
- Golang's `big.Int` has some operations which are more costly than others. For example, doing `n.Exp(n, two, P)` is more costly than doing `n.Mul(n, n); modCoordinate(n)`. This holds for squaring and cubing, but exponents beyond 3 require `.Exp` for the best performance.

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
