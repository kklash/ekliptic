package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestEqualJacobi(t *testing.T) {
	for i, vector1 := range test_vectors.JacobiPointVectors {
		for j, vector2 := range test_vectors.JacobiPointVectors {
			match := EqualJacobi(
				vector1.JacobiX, vector1.JacobiY, vector1.JacobiZ,
				vector2.JacobiX, vector2.JacobiY, vector2.JacobiZ,
			)

			if !match {
				t.Errorf("jacobi-equality check failed for points %d:%d", i, j)
			}
		}
	}
}

func TestEqualAffine(t *testing.T) {
	x1 := new(big.Int)
	y1 := new(big.Int)
	z1 := new(big.Int)
	x2 := new(big.Int)
	y2 := new(big.Int)
	z2 := new(big.Int)

	for i, vector1 := range test_vectors.JacobiPointVectors {
		for j, vector2 := range test_vectors.JacobiPointVectors {
			x1.Set(vector1.JacobiX)
			y1.Set(vector1.JacobiY)
			z1.Set(vector1.JacobiZ)
			x2.Set(vector2.JacobiX)
			y2.Set(vector2.JacobiY)
			z2.Set(vector2.JacobiZ)

			ToAffine(x1, y1, z1)
			ToAffine(x2, y2, z2)

			if !EqualAffine(x1, y1, x2, y2) {
				t.Errorf("affine-equality check failed for points %d:%d", i, j)
			}
		}
	}
}

func BenchmarkEqualJacobi(b *testing.B) {
	vector1 := test_vectors.JacobiPointVectors[0]
	vector2 := test_vectors.JacobiPointVectors[1]

	for i := 0; i < b.N; i++ {
		EqualJacobi(
			vector1.JacobiX, vector1.JacobiY, vector1.JacobiZ,
			vector2.JacobiX, vector2.JacobiY, vector2.JacobiZ,
		)
	}
}

func BenchmarkEqualAffine(b *testing.B) {
	vector1 := test_vectors.JacobiPointVectors[0]
	vector2 := test_vectors.JacobiPointVectors[1]

	x1 := new(big.Int).Set(vector1.JacobiX)
	y1 := new(big.Int).Set(vector1.JacobiY)
	z1 := new(big.Int).Set(vector1.JacobiZ)
	x2 := new(big.Int).Set(vector2.JacobiX)
	y2 := new(big.Int).Set(vector2.JacobiY)
	z2 := new(big.Int).Set(vector2.JacobiZ)
	ToAffine(x1, y1, z1)
	ToAffine(x2, y2, z2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		EqualAffine(x1, y1, x2, y2)
	}
}

// Benchmarks performance of comparing points by converting to affine.
// For comparison against EqualJacobi.
func BenchmarkEqualAffine_WithConversion(b *testing.B) {
	vector1 := test_vectors.JacobiPointVectors[0]
	vector2 := test_vectors.JacobiPointVectors[1]

	x1 := new(big.Int).Set(vector1.JacobiX)
	y1 := new(big.Int).Set(vector1.JacobiY)
	z1 := new(big.Int).Set(vector1.JacobiZ)
	x2 := new(big.Int).Set(vector2.JacobiX)
	y2 := new(big.Int).Set(vector2.JacobiY)
	z2 := new(big.Int).Set(vector2.JacobiZ)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		x1.Set(vector1.JacobiX)
		y1.Set(vector1.JacobiY)
		z1.Set(vector1.JacobiZ)
		x2.Set(vector2.JacobiX)
		y2.Set(vector2.JacobiY)
		z2.Set(vector2.JacobiZ)

		ToAffine(x1, y1, z1)
		ToAffine(x2, y2, z2)
		EqualAffine(x1, y1, x2, y2)
	}
}
