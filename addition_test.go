package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestAddJacobi(t *testing.T) {
	for i, vector := range test_vectors.JacobiAdditionVectors {
		x3 := new(big.Int)
		y3 := new(big.Int)
		z3 := new(big.Int)

		AddJacobi(
			vector.X1, vector.Y1, vector.Z1,
			vector.X2, vector.Y2, vector.Z2,
			x3, y3, z3,
		)

		if !equal(x3, vector.X3) || !equal(y3, vector.Y3) || !equal(z3, vector.Z3) {
			t.Errorf(`jacobi point addition failed for vector %d - Got:
	x3: %x
	y3: %x
	z3: %x
Wanted:
	x3: %x
	y3: %x
	z3: %x
`, i, x3, y3, z3, vector.X3, vector.Y3, vector.Z3)
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

		x3 := new(big.Int)
		y3 := new(big.Int)
		AddAffine(x1, y1, x2, y2, x3, y3)

		if !EqualAffine(x3, y3, expectedX, expectedY) {
			t.Errorf(`affine point addition failed for vector %d - Got:
	x3: %x
	y3: %x
Wanted:
	x3: %x
	y3: %x
`, i, x3, y3, expectedX, expectedY)
		}
	}
}

// Benchmarks adding jacobi points where z = 1 for both points.
func BenchmarkAddJacobi_Z1(b *testing.B) {
	x3 := new(big.Int)
	y3 := new(big.Int)
	z3 := new(big.Int)

	// where both points' z = 1
	vector := test_vectors.JacobiAdditionVectors[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddJacobi(
			vector.X1, vector.Y1, vector.Z1,
			vector.X2, vector.Y2, vector.Z2,
			x3, y3, z3,
		)
	}
}

// Benchmarks adding jacobi points where z > 1 for both points.
func BenchmarkAddJacobi_LargeZ(b *testing.B) {
	x3 := new(big.Int)
	y3 := new(big.Int)
	z3 := new(big.Int)

	// where both points' z > 1
	vector := test_vectors.JacobiAdditionVectors[118]

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		AddJacobi(
			vector.X1, vector.Y1, vector.Z1,
			vector.X2, vector.Y2, vector.Z2,
			x3, y3, z3,
		)
	}
}

func addJacobiWithAllocate(
	x1, y1, z1 *big.Int,
	x2, y2, z2 *big.Int,
) (x3, y3, z3 *big.Int) {
	x3 = new(big.Int)
	y3 = new(big.Int)
	z3 = new(big.Int)

	AddJacobi(
		x1, y1, z1,
		x2, y2, z2,
		x3, y3, z3,
	)

	return x3, y3, z3
}

// Benchmarks adding jacobi point, allocating return values each time
func BenchmarkAddJacobi_WithAllocate(b *testing.B) {
	vector := test_vectors.JacobiAdditionVectors[118] // both points' z > 1

	for i := 0; i < b.N; i++ {
		addJacobiWithAllocate(
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

	x3 := new(big.Int)
	y3 := new(big.Int)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddAffine(x1, y1, x2, y2, x3, y3)
	}
}
