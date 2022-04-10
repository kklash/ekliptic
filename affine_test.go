package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestToAffine(t *testing.T) {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	for i, vector := range test_vectors.JacobiPointVectors {
		x.Set(vector.JacobiX)
		y.Set(vector.JacobiY)
		z.Set(vector.JacobiZ)

		ToAffine(x, y, z)

		if !equal(x, vector.X) || !equal(y, vector.Y) {
			t.Errorf(`jacobian to affine point conversion failed for vector %d. Got:
	x: %.64x
	y: %.64x
Wanted:
	x: %.64x
	y: %.64x
`, i, x, y, vector.X, vector.Y)
		}
	}
}

func BenchmarkToAffine(b *testing.B) {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)

	vector := test_vectors.JacobiPointVectors[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Set(vector.JacobiX)
		y.Set(vector.JacobiY)
		z.Set(vector.JacobiZ)

		ToAffine(x, y, z)
	}
}
