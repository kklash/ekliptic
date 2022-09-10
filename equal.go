package ekliptic

import "math/big"

func equal(v1, v2 *big.Int) bool {
	if v1 == nil || v2 == nil {
		return v1 == v2
	}

	return v1.Cmp(v2) == 0
}

// EqualJacobi tests whether two Jacobian points are equivalent to the same affine point,
// without the performance penalty of actually converting both points to affine format.
// It returns true if:
//
//	x2 * z1² == x1 * z2²
//	y2 * z1³ == y1 * z2³
//
// Affine coordinates are calculated as:
//
//	Ax1 = x1 / z1²
//	Ax2 = x2 / z2²
//	Ay1 = y1 / z1³
//	Ay2 = y2 / z2³
//
// Thus, we can re-arrange operations to avoid costly division,
// and determine if the x-coordinate matches:
//
//	Ax1 = Ax2
//	x2 / z2² = x1 / z1²
//	(x2 * z1²) / z2² = x1
//	x2 * z1² = x1 * z2²
//
// Same for the y-coordinate:
//
//	Ay1 = Ay2
//	(y2 * z1³) / z2³ = y1
//	y2 * z1³ = y1 * z2³
//
// This approach provides a 5x speedup compared to affine conversion.
func EqualJacobi(
	x1, y1, z1 *big.Int,
	x2, y2, z2 *big.Int,
) bool {
	if equal(z1, z2) {
		return equal(x1, x2) && equal(y1, y2)
	}

	// z1² and z2²
	z1_pow2 := new(big.Int).Mul(z1, z1)
	z2_pow2 := new(big.Int).Mul(z2, z2)

	// u1 = x1 * z2²
	u1 := new(big.Int).Mul(x1, z2_pow2)
	modCoordinate(u1)

	// u2 = x2 * z1²
	u2 := new(big.Int).Mul(x2, z1_pow2)
	modCoordinate(u2)

	// Ax1 != Ax2
	if !equal(u1, u2) {
		return false
	}

	// z1³
	z1_pow3 := z1_pow2.Mul(z1_pow2, z1)
	z1_pow2 = nil

	// z2³
	z2_pow3 := z2_pow2.Mul(z2_pow2, z2)
	z2_pow2 = nil

	// s1 = y1 * z2³
	s1 := z2_pow3.Mul(y1, z2_pow3)
	z2_pow3 = nil
	modCoordinate(s1)

	// s2 = y2 * z1³
	s2 := z1_pow3.Mul(y2, z1_pow3)
	z1_pow3 = nil
	modCoordinate(s2)

	return equal(s1, s2)
}

// EqualAffine tests the equality of two affine points.
func EqualAffine(
	x1, y1 *big.Int,
	x2, y2 *big.Int,
) bool {
	return equal(x1, x2) && equal(y1, y2)
}
