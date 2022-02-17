package ekliptic

import "math/big"

// Negate sets the given Y-coordinate to its negative value, modulo P.
// This can be used to negate a point, because:
//
//  (x, y) + (x, -y) = 0
//
// y is expected to be within range [0..P-1].
func Negate(y *big.Int) {
	if !equal(y, zero) {
		y.Sub(Secp256k1_P, y)
	}
}
