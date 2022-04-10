package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestMultiplyJacobi(t *testing.T) {
	resultX := new(big.Int)
	resultY := new(big.Int)
	resultZ := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			resultX, resultY, resultZ,
			nil,
		)
		ToAffine(resultX, resultY, resultZ)

		if !equal(resultX, vector.X2) || !equal(resultY, vector.Y2) {
			t.Errorf(`jacobi multiplication failed for vector %d. Got:
	x: %.64x
	y: %.64x
Wanted:
	x: %.64x
	y: %.64x
`, i, resultX, resultY, vector.X2, vector.Y2)
		}
	}
}

func TestMultiplyJacobi_MemSafety(t *testing.T) {
	x1 := new(big.Int)
	y1 := new(big.Int)
	z1 := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		x1.Set(vector.X1)
		y1.Set(vector.Y1)
		z1.Set(one)

		MultiplyJacobi(
			x1, y1, z1,
			vector.K,
			x1, y1, z1,
			nil,
		)
		ToAffine(x1, y1, z1)

		if !equal(x1, vector.X2) || !equal(y1, vector.Y2) {
			t.Errorf(`jacobi memory-safe multiplication failed for vector %d. Got:
	x: %.64x
	y: %.64x
Wanted:
	x: %.64x
	y: %.64x
`, i, x1, y1, vector.X2, vector.Y2)
		}
	}
}

func TestMultiplyAffine(t *testing.T) {
	resultX := new(big.Int)
	resultY := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		MultiplyAffine(
			vector.X1, vector.Y1,
			vector.K,
			resultX, resultY,
			nil,
		)

		if !equal(resultX, vector.X2) || !equal(resultY, vector.Y2) {
			t.Errorf(`affine multiplication failed for vector %d. Got:
	x: %.64x
	y: %.64x
Wanted:
	x: %.64x
	y: %.64x
`, i, resultX, resultY, vector.X2, vector.Y2)
		}
	}
}

func TestMultiplyAffineNaive(t *testing.T) {
	resultX := new(big.Int)
	resultY := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		MultiplyAffineNaive(
			vector.X1, vector.Y1,
			vector.K,
			resultX, resultY,
			nil,
		)

		if !equal(resultX, vector.X2) || !equal(resultY, vector.Y2) {
			t.Errorf(`affine naive multiplication failed for vector %d. Got:
	x: %.64x
	y: %.64x
Wanted:
	x: %.64x
	y: %.64x
`, i, resultX, resultY, vector.X2, vector.Y2)
		}
	}
}

func BenchmarkMultiplyJacobi(b *testing.B) {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)

	vector := test_vectors.AffineMultiplicationVectors[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			x, y, z,
			nil,
		)
	}
}

func BenchmarkMultiplyJacobi_Precomputed(b *testing.B) {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)

	vector := test_vectors.AffineMultiplicationVectors[0]

	precomputes := ComputePointDoubles(vector.X1, vector.Y1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			x, y, z,
			precomputes,
		)
	}
}

func BenchmarkMultiplyAffine(b *testing.B) {
	x := new(big.Int)
	y := new(big.Int)

	vector := test_vectors.AffineMultiplicationVectors[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyAffine(
			vector.X1, vector.Y1,
			vector.K,
			x, y,
			nil,
		)
	}
}

func BenchmarkMultiplyAffineNaive(b *testing.B) {
	x := new(big.Int)
	y := new(big.Int)

	vector := test_vectors.AffineMultiplicationVectors[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyAffineNaive(
			vector.X1, vector.Y1,
			vector.K,
			x, y,
			nil,
		)
	}
}
