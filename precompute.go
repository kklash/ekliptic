package ekliptic

import "math/big"

//go:generate go run ./genprecompute -o precomputed_doubles.go

// PrecomputedDoubles is a slice of points which are precomputed by doubling a given point.
// You can pass a set of PrecomputedDoubles to MultiplyJacobi to speed up multiplication operations.
type PrecomputedDoubles = [][2]*big.Int

// ComputePointDoubles is used to precompute the results of doubling a given secp256k1 point. The
// computed results can be passed to MultiplyJacobi to speed up multiplication operations on that same point.
func ComputePointDoubles(genX, genY *big.Int) PrecomputedDoubles {
	precomputedBasePointDoubles := make(PrecomputedDoubles, 256)

	x := new(big.Int).Set(genX)
	y := new(big.Int).Set(genY)
	z := new(big.Int).Set(one)

	for i := 0; i < len(precomputedBasePointDoubles); i++ {
		precomputedBasePointDoubles[i] = [2]*big.Int{
			new(big.Int).Set(x),
			new(big.Int).Set(y),
		}
		x, y, z = DoubleJacobi(x, y, z)
		ToAffine(x, y, z)
	}

	return precomputedBasePointDoubles
}
