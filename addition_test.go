package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestAddJacobi(t *testing.T) {
	for i, vector := range test_vectors.JacobiAdditionVectors {
		x3, y3, z3 := AddJacobi(
			vector.X1, vector.Y1, vector.Z1,
			vector.X2, vector.Y2, vector.Z2,
		)

		if !equal(x3, vector.X3) || !equal(y3, vector.Y3) || !equal(z3, vector.Z3) {
			t.Errorf(`jacobi point addition failed for vector %d - Got:
	x3: %.64x
	y3: %.64x
	z3: %.64x
Wanted:
	x3: %.64x
	y3: %.64x
	z3: %.64x
`, i, x3, y3, z3, vector.X3, vector.Y3, vector.Z3)
		}
	}
}

func TestAddJacobi_MemSafety(t *testing.T) {
	for i, vector := range test_vectors.JacobiAdditionVectors {
		originalX1 := new(big.Int).Set(vector.X1)
		originalY1 := new(big.Int).Set(vector.Y1)
		originalZ1 := new(big.Int).Set(vector.Z1)

		originalX2 := new(big.Int).Set(vector.X2)
		originalY2 := new(big.Int).Set(vector.Y2)
		originalZ2 := new(big.Int).Set(vector.Z2)

		AddJacobi(
			vector.X1, vector.Y1, vector.Z1,
			vector.X2, vector.Y2, vector.Z2,
		)

		if !equal(vector.X1, originalX1) || !equal(vector.Y1, originalY1) || !equal(vector.Z1, originalZ1) {
			t.Errorf(`jacobi memory-safe point addition failed for vector %d - Got:
	x1: %.64x
	y1: %.64x
	z1: %.64x
Wanted:
	x1: %.64x
	y1: %.64x
	z1: %.64x
`, i, vector.X1, vector.Y1, vector.Z1, originalX1, originalY1, originalZ1)
		}

		if !equal(vector.X2, originalX2) || !equal(vector.Y2, originalY2) || !equal(vector.Z2, originalZ2) {
			t.Errorf(`jacobi memory-safe point addition failed for vector %d - Got:
	x2: %.64x
	y2: %.64x
	z2: %.64x
Wanted:
	x2: %.64x
	y2: %.64x
	z2: %.64x
`, i, vector.X2, vector.Y2, vector.Z2, originalX2, originalY2, originalZ2)
		}
	}
}

func TestSubJacobi(t *testing.T) {
	for i, vector := range test_vectors.JacobiAdditionVectors {
		x1, y1, z1 := SubJacobi(
			vector.X3, vector.Y3, vector.Z3,
			vector.X2, vector.Y2, vector.Z2,
		)

		if !EqualJacobi(x1, y1, z1, vector.X1, vector.Y1, vector.Z1) {
			ToAffine(x1, y1, z1)

			expectedX1 := new(big.Int).Set(vector.X1)
			expectedY1 := new(big.Int).Set(vector.Y1)
			expectedZ1 := new(big.Int).Set(vector.Z1)
			ToAffine(expectedX1, expectedY1, expectedZ1)

			t.Errorf(`jacobi point subtraction failed for vector %d - Got affine points:
	x1: %.64x
	y1: %.64x
Wanted:
	x1: %.64x
	y1: %.64x
`, i, x1, y1, expectedX1, expectedY1)
		}
	}
}

func TestAddAffine(t *testing.T) {
	for i, vector := range test_vectors.JacobiAdditionVectors {
		x1 := new(big.Int).Set(vector.X1)
		y1 := new(big.Int).Set(vector.Y1)
		z1 := new(big.Int).Set(vector.Z1)
		x2 := new(big.Int).Set(vector.X2)
		y2 := new(big.Int).Set(vector.Y2)
		z2 := new(big.Int).Set(vector.Z2)
		expectedX := new(big.Int).Set(vector.X3)
		expectedY := new(big.Int).Set(vector.Y3)
		expectedZ := new(big.Int).Set(vector.Z3)

		ToAffine(x1, y1, z1)
		ToAffine(x2, y2, z2)
		ToAffine(expectedX, expectedY, expectedZ)

		x3, y3 := AddAffine(x1, y1, x2, y2)

		if !EqualAffine(x3, y3, expectedX, expectedY) {
			t.Errorf(`affine point addition failed for vector %d - Got:
	x3: %.64x
	y3: %.64x
Wanted:
	x3: %.64x
	y3: %.64x
`, i, x3, y3, expectedX, expectedY)
		}
	}
}

func TestAddAffine_MemSafety(t *testing.T) {
	// Test memory safety when result pointers are also input parameters.
	for i, vector := range test_vectors.JacobiAdditionVectors {
		x1 := new(big.Int).Set(vector.X1)
		y1 := new(big.Int).Set(vector.Y1)
		z1 := new(big.Int).Set(vector.Z1)
		ToAffine(x1, y1, z1)

		x2 := new(big.Int).Set(vector.X2)
		y2 := new(big.Int).Set(vector.Y2)
		z2 := new(big.Int).Set(vector.Z2)
		ToAffine(x2, y2, z2)

		originalX1 := new(big.Int).Set(x1)
		originalY1 := new(big.Int).Set(y1)

		originalX2 := new(big.Int).Set(x2)
		originalY2 := new(big.Int).Set(y2)

		AddAffine(
			x1, y1,
			x2, y2,
		)

		if !EqualAffine(x1, y1, originalX1, originalY1) {
			t.Errorf(`affine memory-safe point addition failed for vector %d - Got:
	x1: %.64x
	y1: %.64x
Wanted:
	x1: %.64x
	y1: %.64x
`, i, x1, y1, originalX1, originalY1)
		}

		if !EqualAffine(x2, y2, originalX2, originalY2) {
			t.Errorf(`affine memory-safe point addition failed for vector %d - Got:
	x2: %.64x
	y2: %.64x
Wanted:
	x2: %.64x
	y2: %.64x
`, i, x2, y2, originalX2, originalY2)
		}
	}
}

func TestSubAffine(t *testing.T) {
	for i, vector := range test_vectors.JacobiAdditionVectors {
		x3 := new(big.Int).Set(vector.X3)
		y3 := new(big.Int).Set(vector.Y3)
		z3 := new(big.Int).Set(vector.Z3)
		x2 := new(big.Int).Set(vector.X2)
		y2 := new(big.Int).Set(vector.Y2)
		z2 := new(big.Int).Set(vector.Z2)
		expectedX1 := new(big.Int).Set(vector.X1)
		expectedY1 := new(big.Int).Set(vector.Y1)
		expectedZ1 := new(big.Int).Set(vector.Z1)

		ToAffine(x3, y3, z3)
		ToAffine(x2, y2, z2)
		ToAffine(expectedX1, expectedY1, expectedZ1)

		x1, y1 := SubAffine(x3, y3, x2, y2)

		if !EqualAffine(x1, y1, expectedX1, expectedY1) {
			t.Errorf(`affine point addition failed for vector %d - Got:
	x3: %.64x
	y3: %.64x
Wanted:
	x3: %.64x
	y3: %.64x
`, i, x1, y1, expectedX1, expectedY1)
		}
	}
}

// Benchmarks adding jacobi points where z = 1 for both points.
func BenchmarkAddJacobi_Z1(b *testing.B) {
	// where both points' z = 1
	vector := test_vectors.JacobiAdditionVectors[0]

	for i := 0; i < b.N; i++ {
		AddJacobi(
			vector.X1, vector.Y1, vector.Z1,
			vector.X2, vector.Y2, vector.Z2,
		)
	}
}

// Benchmarks adding jacobi points where z > 1 for both points.
func BenchmarkAddJacobi_LargeZ(b *testing.B) {
	// where both points' z > 1
	vector := test_vectors.JacobiAdditionVectors[118]

	for i := 0; i < b.N; i++ {
		AddJacobi(
			vector.X1, vector.Y1, vector.Z1,
			vector.X2, vector.Y2, vector.Z2,
		)
	}
}

func BenchmarkAddAffine(b *testing.B) {
	vector := test_vectors.JacobiAdditionVectors[0]

	x1 := new(big.Int).Set(vector.X1)
	y1 := new(big.Int).Set(vector.Y1)
	z1 := new(big.Int).Set(vector.Z1)

	x2 := new(big.Int).Set(vector.X2)
	y2 := new(big.Int).Set(vector.Y2)
	z2 := new(big.Int).Set(vector.Z2)

	ToAffine(x1, y1, z1)
	ToAffine(x2, y2, z2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddAffine(x1, y1, x2, y2)
	}
}
