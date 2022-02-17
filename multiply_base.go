package ekliptic

import (
	"math/big"
)

var basePointPrecomputations PrecomputedDoubles

// MultiplyBasePoint multiplies the secp256k1 generator base point by the
// given integer k, and sets x and y to resulting affine point. This uses
// precomputed doubles for the secp256k1 base point to speed up multiplications.
//
// This function is used to derive the public key for a private key k, among other uses.
func MultiplyBasePoint(k, x, y *big.Int) {
	MultiplyAffine(
		Secp256k1_GeneratorX, Secp256k1_GeneratorY,
		k,
		x, y,
		basePointPrecomputations,
	)
}
