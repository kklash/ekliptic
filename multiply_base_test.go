package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestMultiplyBasePoint(t *testing.T) {
	x := new(big.Int)
	y := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		if !equal(vector.X1, Secp256k1_GeneratorX) ||
			!equal(vector.Y1, Secp256k1_GeneratorY) {
			continue
		}

		MultiplyBasePoint(vector.K, x, y)

		if !equal(x, vector.X2) || !equal(y, vector.Y2) {
			t.Errorf(`multiplying base point failed for vector %d. Got
	x: %x
	y: %x
Wanted:
	x: %x
	y: %x
	`, i, x, y, vector.X2, vector.Y2)
		}
	}
}

func BenchmarkMultiplyBasePoint(b *testing.B) {
	vector := test_vectors.AffineMultiplicationVectors[0]
	x := new(big.Int)
	y := new(big.Int)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyBasePoint(vector.K, x, y)
	}
}
