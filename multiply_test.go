package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestMultiplyJacobi(t *testing.T) {
	for i, vector := range test_vectors.AffineMultiplicationVectors {
		resultX, resultY, resultZ := MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
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

func TestMultiplyJacobi_Precomputed(t *testing.T) {
	for i, vector := range test_vectors.AffineMultiplicationVectors {
		table := NewPrecomputedTable(vector.X1, vector.Y1)
		resultX, resultY, resultZ := MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			table,
		)
		ToAffine(resultX, resultY, resultZ)

		if !equal(resultX, vector.X2) || !equal(resultY, vector.Y2) {
			t.Errorf(`jacobi precomputed multiplication failed for vector %d. Got:
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
	// inputs and outputs are the same pointers
	for i, vector := range test_vectors.AffineMultiplicationVectors {
		originalX1 := new(big.Int).Set(vector.X1)
		originalY1 := new(big.Int).Set(vector.Y1)

		originalK := new(big.Int).Set(vector.K)

		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			nil,
		)

		if !equal(vector.X1, originalX1) || !equal(vector.Y1, originalY1) {
			t.Errorf(`jacobi memory-safe multiplication failed for vector %d. Got:
       x1: %.64x
       y1: %.64x
Wanted:
       x1: %.64x
       x1: %.64x
`, i, vector.X1, vector.Y1, originalX1, originalY1)
		}
		if !equal(vector.K, originalK) {
			t.Errorf(`jacobi memory safe multiplication failed for vector %d. Got:
			 k: %.64x
Wanted:
			 k: %.64x
`, i, vector.K, originalK)
		}
	}
}

func TestMultiplyAffine(t *testing.T) {
	for i, vector := range test_vectors.AffineMultiplicationVectors {
		resultX, resultY := MultiplyAffine(
			vector.X1, vector.Y1,
			vector.K,
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
	for i, vector := range test_vectors.AffineMultiplicationVectors {
		resultX, resultY := MultiplyAffineNaive(
			vector.X1, vector.Y1,
			vector.K,
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
	vector := test_vectors.AffineMultiplicationVectors[0]

	for i := 0; i < b.N; i++ {
		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			nil,
		)
	}
}

func BenchmarkMultiplyJacobi_Precomputed(b *testing.B) {
	vector := test_vectors.AffineMultiplicationVectors[0]

	precomputes := NewPrecomputedTable(vector.X1, vector.Y1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			precomputes,
		)
	}
}

func BenchmarkMultiplyAffine(b *testing.B) {
	vector := test_vectors.AffineMultiplicationVectors[0]

	for i := 0; i < b.N; i++ {
		MultiplyAffine(
			vector.X1, vector.Y1,
			vector.K,
			nil,
		)
	}
}

func BenchmarkMultiplyAffineNaive(b *testing.B) {
	vector := test_vectors.AffineMultiplicationVectors[0]

	for i := 0; i < b.N; i++ {
		MultiplyAffineNaive(
			vector.X1, vector.Y1,
			vector.K,
		)
	}
}
