package ekliptic

import "math/big"

// IsOnCurveAffine determines if the affine point at (x, y) is a valid point on the secp256k1 curve.
// It checks for equality using the Weierstrass curve equation:
//
//  y² = x³ + ax + b
func IsOnCurveAffine(x, y *big.Int) bool {
	if equal(x, zero) && equal(y, zero) {
		return true
	}

	// y²
	left := new(big.Int).Mul(y, y)
	mod(left)

	// x³ + ax + b
	right := new(big.Int).Mul(x, x)
	right.Mul(right, x)
	right.Add(right, Secp256k1_B)
	mod(right)

	return equal(left, right)
}

// IsOnCurveJacobi determines if the jacobian point at (x, y, z) is a valid point on the secp256k1 curve.
// If z is nil, it is assumed to be 1, meaning x and y are affine coordinates. It uses the Weierstrass
// curve equation to determine whether the given point is valid:
//
//  y² = x³ + ax + b
//
// Since a = 0 in secp256k1, we can simplify this to
//
//  y² = x³ + b
//
// When using Jacobian coordinates, x and y are expressed as ratios with respect to z:
//
//  x = x / z²
//  y = y / z³
//
// Substituting and simplifying, we arrive at a much cleaner equation we can easily solve:
//
//  (y/z³)² = (x/z²)³ + b
//  y²/z⁶ = x³/z⁶ + b
//  y² = x³ + z⁶b
//
// This approach is 2.5x more performant than converting a Jacobian point to affine
// coordinates first, because it does not require modular inversion to arrive at a result.
func IsOnCurveJacobi(x, y, z *big.Int) bool {
	if equal(x, zero) && equal(y, zero) {
		return true
	}
	if z == nil || equal(z, one) {
		return IsOnCurveAffine(x, y)
	}

	// y²
	left := new(big.Int).Mul(y, y)
	mod(left)

	// z⁶b
	// This is more efficient for larger powers than repeated .Mul() calls
	z6b := new(big.Int).Exp(z, six, Secp256k1_P)
	z6b.Mul(z6b, Secp256k1_B)

	// x³ + z⁶b
	right := new(big.Int).Mul(x, x)
	right.Mul(right, x)
	right.Add(right, z6b)
	mod(right)

	return equal(left, right)
}
