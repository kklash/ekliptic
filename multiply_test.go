package ekliptic

import (
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

	precomputes := ComputePointDoubles(vector.X1, vector.Y1)

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
			nil,
		)
	}
}
