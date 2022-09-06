package ekliptic

import (
	"math/big"
)

// MultiplyJacobi multiplies the given Jacobian point (x1, y1, z1) by the scalar value k
// in constant time.
//
// It returns the resulting Jacobian point (x2, y2, z2).
//
// Callers can construct and pass a PrecomputedTable which massively boosts performance
// of a MultiplyJacobi call, at the cost of a larger up-front computational investment to
// build the precomputations. If a you plan to multiply the same point several times,
// precomputing is definitely worthwhile.
//
// If a PrecomputedTable is passed, MultiplyJacobi will use the windowed multiplication
// method for fast computation. Otherwise, it will use the Montgomery Ladder algorithm.
//
// MultiplyJacobi checks and panics if the given point you are multiplying is not actually
// on the secp256k1 curve, as this could leak private data about the scalar value k.
func MultiplyJacobi(
	x1, y1, z1 *big.Int,
	k *big.Int,
	precomputedTable PrecomputedTable,
) (x2, y2, z2 *big.Int) {
	if !IsOnCurveJacobi(x1, y1, z1) {
		panic("MultiplyJacobi: refusing to multiply point not on the curve; this could leak private data")
	}

	if precomputedTable != nil {
		return multiplyJacobiTable(x1, y1, z1, k, precomputedTable)
	}

	x2 = new(big.Int)
	y2 = new(big.Int)
	z2 = new(big.Int)

	dummyX := new(big.Int).Set(x1)
	dummyY := new(big.Int).Set(y1)
	dummyZ := new(big.Int).Set(z1)

	for i := 255; i >= 0; i-- {
		if k.Bit(i) > 0 {
			x2, y2, z2 = AddJacobi(
				x2, y2, z2,
				dummyX, dummyY, dummyZ,
			)
			dummyX, dummyY, dummyZ = DoubleJacobi(dummyX, dummyY, dummyZ)
		} else {
			dummyX, dummyY, dummyZ = AddJacobi(
				x2, y2, z2,
				dummyX, dummyY, dummyZ,
			)
			x2, y2, z2 = DoubleJacobi(x2, y2, z2)
		}
	}
	return
}

// multiplyJacobiTable multiplies the given Jacobian point by the scalar value k using
// the windowed multiplication method, with a window size of 4. This function expects
// to receive a precomputed multiplication table for lookups of point doubles and the
// products of point doubles.
func multiplyJacobiTable(
	x1, y1, z1 *big.Int,
	k *big.Int,
	precomputedTable PrecomputedTable,
) (x2, y2, z2 *big.Int) {
	kBytes := k.FillBytes(make([]byte, 32))
	windows := make([]byte, 64)
	for i := range kBytes {
		windows[i*2] = kBytes[i] >> 4
		windows[i*2+1] = kBytes[i] & 0b1111
	}

	x2 = new(big.Int)
	y2 = new(big.Int)
	z2 = new(big.Int)

	dummyX := new(big.Int)
	dummyY := new(big.Int)
	dummyZ := new(big.Int)

	rowLen := len(precomputedTable[0])

	for i, d := range windows {
		row := precomputedTable[63-i]
		if d > 0 {
			point := row[d]
			x2, y2, z2 = AddJacobi(
				x2, y2, z2,
				point[0], point[1], one,
			)
		} else {
			// We add an arbitrary point to the dummy point to ensure constant
			// time operation.
			c := i%(rowLen-1) + 1 // 0 < c < len(row)
			oppositePoint := row[c]
			dummyX, dummyY, dummyZ = AddJacobi(
				dummyX, dummyY, dummyZ,
				oppositePoint[0], oppositePoint[1], one,
			)
		}
	}

	return
}

// MultiplyAffine multiplies the given affine point (x1, y1) by the scalar value k in constant time.
//
// If a PrecomputedTable is passed, MultiplyAffine will use the windowed multiplication
// method for fast computation. Otherwise, it will use the Montgomery Ladder algorithm.
//
// It returns the resulting affine point (x2, y2).
//
// Callers can construct and pass a PrecomputedTable which massively boosts performance
// of a MultiplyAffine call, at the cost of a larger up-front computational investment to
// build the precomputations. If a you plan to multiply the same point several times,
// precomputing is definitely worthwhile.
//
// MultiplyAffine uses MultiplyJacobi under the hood, as it is about 30% faster than performing affine addition.
func MultiplyAffine(
	x1, y1 *big.Int,
	k *big.Int,
	precomputedTable PrecomputedTable,
) (x2, y2 *big.Int) {
	x2, y2, z2 := MultiplyJacobi(
		x1, y1, one,
		k,
		precomputedTable,
	)

	ToAffine(x2, y2, z2)
	return
}

// MultiplyAffineNaive multiplies the given affine point by the scalar value k, using the Montgomery
// Ladder algorithm in constant-time.
//
// This naive implementation uses affine doubling and addition under the hood, which is not desirable.
// It is made available only as a demonstration of how much faster Jacobian math is for elliptic curve
// multiplication operations. You should be using MultiplyAffine instead.
func MultiplyAffineNaive(
	x1, y1 *big.Int,
	k *big.Int,
) (x2, y2 *big.Int) {
	x2 = new(big.Int).Set(zero)
	y2 = new(big.Int).Set(zero)

	dummyX := new(big.Int).Set(x1)
	dummyY := new(big.Int).Set(y1)

	for i := 255; i >= 0; i-- {
		if k.Bit(i) > 0 {
			x2, y2 = AddAffine(
				x2, y2,
				dummyX, dummyY,
			)
			dummyX, dummyY = DoubleAffine(dummyX, dummyY)
		} else {
			dummyX, dummyY = AddAffine(
				x2, y2,
				dummyX, dummyY,
			)
			x2, y2 = DoubleAffine(x2, y2)
		}
	}
	return
}
