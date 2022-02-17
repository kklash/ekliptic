package ekliptic

import "math/big"

// Jacobian coordinates are a three-dimensional representation of a 2d (affine) point, (Ax, Ay), in terms of three variables: (x, y, z)
// such that:
//  Ax = x / z²
//  Ay = y / z³
//
// ToAffine converts the given jacobian coordinates to affine coordinates,
// normalizing x and y, and setting z = 1. This is an expensive operation,
// as it involves modular inversion of z to perform finite field division.
func ToAffine(x, y, z *big.Int) {
	if equal(z, one) {
		return
	} else if equal(z, zero) {
		x.Set(zero)
		y.Set(zero)
		return
	}

	invert(z)

	x.Mul(x, z)
	x.Mul(x, z)
	mod(x)

	y.Mul(y, z)
	y.Mul(y, z)
	y.Mul(y, z)
	mod(y)

	z.Set(one)
}
