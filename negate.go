package ekliptic

import "math/big"

// Negate returns the additive inverse of the given Y-coordinate modulo the
// curve prime modulus P. This can be used to negate a point, because:
//
//	(x, y) + (x, -y) = 0
//
// y is expected to be within range [0, P-1].
func Negate(y *big.Int) *big.Int {
	if equal(y, zero) {
		return new(big.Int).Set(zero)
	}
	return new(big.Int).Sub(Secp256k1_P, y)
}
