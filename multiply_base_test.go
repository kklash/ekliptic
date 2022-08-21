package ekliptic

import (
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestMultiplyBasePoint(t *testing.T) {
	for i, vector := range test_vectors.AffineMultiplicationVectors {
		if !equal(vector.X1, Secp256k1_GeneratorX) ||
			!equal(vector.Y1, Secp256k1_GeneratorY) {
			continue
		}

		x, y := MultiplyBasePoint(vector.K)

		if !equal(x, vector.X2) || !equal(y, vector.Y2) {
			t.Errorf(`multiplying base point failed for vector %d. Got
	x: %.64x
	y: %.64x
Wanted:
	x: %.64x
	y: %.64x
	`, i, x, y, vector.X2, vector.Y2)
		}
	}
}

func BenchmarkMultiplyBasePoint(b *testing.B) {
	vector := test_vectors.AffineMultiplicationVectors[0]

	for i := 0; i < b.N; i++ {
		MultiplyBasePoint(vector.K)
	}
}
