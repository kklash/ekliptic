package ekliptic

import (
	"math/big"
)

// MultiplyJacobi multiplies the given Jacobian point (x1, y1, z1) by the scalar value k, using the double-and-add algorithm:
//
//  P2 = P1 * k
//  P2 = (P1 * 2^0) + (P1 * 2^1) + (P1 * 2^2) + (P1 * 2^3) ...
//
// It stores the resulting Jacobian point in (x2, y2, z2).
//
// You can pass compute and then pass PrecomputedDoubles which massively boosts performance of successive MultiplyJacobi
// calls, at the cost of a larger up-front computational investment. If you plan to multiply the same point more than
// just once or twice, precomputing is definitely worthwhile.
//
// MultiplyJacobi checks and panics if the given point you are multiplying is not actually on the secp256k1 curve,
// as this could leak private data about the scalar value k.
func MultiplyJacobi(
	x1, y1, z1 *big.Int,
	k *big.Int,
	x2, y2, z2 *big.Int,
	precomputedDoubles PrecomputedDoubles,
) {
	if !IsOnCurveJacobi(x1, y1, z1) {
		panic("MultiplyJacobi: refusing to multiply point not on the curve; this could leak private data")
	}

	// save the bits ahead of time in case k points to the same int as one of the outputs.
	kBits := make([]bool, k.BitLen())
	for i := range kBits {
		kBits[i] = k.Bit(i) > 0
	}

	doubleX := new(big.Int).Set(x1)
	doubleY := new(big.Int).Set(y1)
	doubleZ := new(big.Int).Set(z1)

	x2.Set(zero)
	y2.Set(zero)
	z2.Set(zero)

	for i, shouldAdd := range kBits {
		if shouldAdd {
			AddJacobi(
				x2, y2, z2,
				doubleX, doubleY, doubleZ,
				x2, y2, z2,
			)
		}

		if i+1 < len(precomputedDoubles) {
			doubleX.Set(precomputedDoubles[i+1][0])
			doubleY.Set(precomputedDoubles[i+1][1])
			doubleZ.Set(one)
		} else {
			DoubleJacobi(
				doubleX, doubleY, doubleZ,
				doubleX, doubleY, doubleZ,
			)
		}
	}
}

// MultiplyAffine multiplies the given affine point by the scalar value k, using the double-and-add algorithm:
//
//  P2 = P1 * k
//  P2 = (P1 * 2^0) + (P1 * 2^1) + (P1 * 2^2) + (P1 * 2^3) ...
//
// You can pass it a PrecomputedDoubles struct which massively boosts performance at the cost of
// a larger up-front computational investment. If you plan to multiply the same point more than
// just once or twice, precomputing is definitely worthwhile.
//
// MultiplyAffine uses MultiplyJacobi under the hood, as it is about 30% faster than performing affine addition.
func MultiplyAffine(
	x1, y1 *big.Int,
	k *big.Int,
	x2, y2 *big.Int,
	precomputedDoubles PrecomputedDoubles,
) {
	z := new(big.Int).Set(one)

	MultiplyJacobi(
		x1, y1, z,
		k,
		x2, y2, z,
		precomputedDoubles,
	)

	ToAffine(x2, y2, z)
}

// MultiplyAffineNaive multiplies the given affine point by the scalar value k, using the double-and-add algorithm.
//
// This naive implementation uses affine doubling and addition under the hood, which is not desirable.
// It is made available only as a demonstration of how much faster Jacobian math is for elliptic curve
// multiplication operations. You should be using MultiplyAffine instead.
func MultiplyAffineNaive(
	x1, y1 *big.Int,
	k *big.Int,
	x2, y2 *big.Int,
	precomputedDoubles PrecomputedDoubles,
) {
	doubleX := new(big.Int).Set(x1)
	doubleY := new(big.Int).Set(y1)

	x2.Set(zero)
	y2.Set(zero)

	for i := 0; i < k.BitLen(); i++ {
		if k.Bit(i) > 0 {
			AddAffine(
				x2, y2,
				doubleX, doubleY,
				x2, y2,
			)
		}

		if i+1 < len(precomputedDoubles) {
			doubleX.Set(precomputedDoubles[i+1][0])
			doubleY.Set(precomputedDoubles[i+1][1])
		} else {
			DoubleAffine(
				doubleX, doubleY,
				doubleX, doubleY,
			)
		}
	}

}
