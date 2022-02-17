// Package ekliptic provides primitives for cryptographic operations on the secp256k1 curve,
// with zero dependencies and excellent performance. It provides both Affine and Jacobian
// interfaces for elliptic curve operations. Aims to facilitate performant low-level operations
// on secp256k1 without overengineering or kitchen-sink syndrome.
package ekliptic

import (
	"math/big"
)

var (
	zero  = big.NewInt(0)
	one   = big.NewInt(1)
	two   = big.NewInt(2)
	three = big.NewInt(3)
	four  = big.NewInt(4)
	five  = big.NewInt(5)
	six   = big.NewInt(6)
	seven = big.NewInt(7)
	eight = big.NewInt(8)
)

var (
	// The parameters of the secp256k1 elliptic curve.
	Secp256k1_B          = seven
	Secp256k1_P          = hexint("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
	Secp256k1_CurveOrder = hexint("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")

	// Secp256k1_GeneratorX and Secp256k1_GeneratorY descirbe the secp256k1
	// generator point, used for deriving public keys from private keys.
	Secp256k1_GeneratorX = hexint("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798")
	Secp256k1_GeneratorY = hexint("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")

	// Secp256k1_CurveOrderHalf is half of Secp256k1_CurveOrder (rounded down).
	Secp256k1_CurveOrderHalf = new(big.Int).Rsh(Secp256k1_CurveOrder, 1)

	// (P + 1) / 4. Any scalar value raised to this power modulo P will have been square rooted
	// within the finite field. https://en.wikipedia.org/wiki/Quadratic_residue#Prime_or_prime_power_modulus
	squareRootExp = hexint("3FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFBFFFFF0C")
)
