package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestDoubleJacobi(t *testing.T) {
	for i, vector := range test_vectors.JacobiDoublingVectors {
		x3 := new(big.Int)
		y3 := new(big.Int)
		z3 := new(big.Int)

		DoubleJacobi(
			vector.X1, vector.Y1, vector.Z1,
			x3, y3, z3,
		)

		if !equal(x3, vector.X3) || !equal(y3, vector.Y3) || !equal(z3, vector.Z3) {
			t.Errorf(`jacobi point doubling failed for vector %d - Got:
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

func TestDoubleJacobi_MemSafety(t *testing.T) {
	for i, vector := range test_vectors.JacobiDoublingVectors {
		x1 := new(big.Int).Set(vector.X1)
		y1 := new(big.Int).Set(vector.Y1)
		z1 := new(big.Int).Set(vector.Z1)

		DoubleJacobi(
			x1, y1, z1,
			x1, y1, z1,
		)

		if !equal(x1, vector.X3) || !equal(y1, vector.Y3) || !equal(z1, vector.Z3) {
			t.Errorf(`jacobi memory-safe point doubling failed for vector %d - Got:
	x3: %.64x
	y3: %.64x
	z3: %.64x
Wanted:
	x3: %.64x
	y3: %.64x
	z3: %.64x
`, i, x1, y1, z1, vector.X3, vector.Y3, vector.Z3)
		}
	}
}

func TestDoubleAffine(t *testing.T) {
	for i, vector := range test_vectors.JacobiDoublingVectors {
		x1 := new(big.Int).Set(vector.X1)
		y1 := new(big.Int).Set(vector.Y1)
		z1 := new(big.Int).Set(vector.Z1)
		expectedX := new(big.Int).Set(vector.X3)
		expectedY := new(big.Int).Set(vector.Y3)
		expectedZ := new(big.Int).Set(vector.Z3)

		ToAffine(x1, y1, z1)
		ToAffine(expectedX, expectedY, expectedZ)

		x3 := new(big.Int)
		y3 := new(big.Int)
		DoubleAffine(x1, y1, x3, y3)

		if !EqualAffine(x3, y3, expectedX, expectedY) {
			t.Errorf(`affine point doubling failed for vector %d - Got:
	x3: %.64x
	y3: %.64x
Wanted:
	x3: %.64x
	y3: %.64x
`, i, x3, y3, expectedX, expectedY)
		}
	}
}

func TestDoubleAffine_MemSafety(t *testing.T) {
	for i, vector := range test_vectors.JacobiDoublingVectors {
		x1 := new(big.Int).Set(vector.X1)
		y1 := new(big.Int).Set(vector.Y1)
		z1 := new(big.Int).Set(vector.Z1)
		expectedX := new(big.Int).Set(vector.X3)
		expectedY := new(big.Int).Set(vector.Y3)
		expectedZ := new(big.Int).Set(vector.Z3)

		ToAffine(x1, y1, z1)
		ToAffine(expectedX, expectedY, expectedZ)

		DoubleAffine(x1, y1, x1, y1)

		if !EqualAffine(x1, y1, expectedX, expectedY) {
			t.Errorf(`affine memory-safe point doubling failed for vector %d - Got:
	x3: %.64x
	y3: %.64x
Wanted:
	x3: %.64x
	y3: %.64x
`, i, x1, y1, expectedX, expectedY)
		}
	}
}

// Benchmarks doubling jacobi points where z = 1
func BenchmarkDoubleJacobi_Z1(b *testing.B) {
	x3 := new(big.Int)
	y3 := new(big.Int)
	z3 := new(big.Int)

	vector := test_vectors.JacobiDoublingVectors[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DoubleJacobi(
			vector.X1, vector.Y1, vector.Z1,
			x3, y3, z3,
		)
	}
}

// Benchmarks doubling jacobi points where z > 1
func BenchmarkDoubleJacobi_LargeZ(b *testing.B) {
	x3 := new(big.Int)
	y3 := new(big.Int)
	z3 := new(big.Int)

	// point to double where z > 1
	vector := test_vectors.JacobiDoublingVectors[2]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DoubleJacobi(
			vector.X1, vector.Y1, vector.Z1,
			x3, y3, z3,
		)
	}
}

func doubleJacobiWithAllocate(
	x1, y1, z1 *big.Int,
) (x3, y3, z3 *big.Int) {
	x3 = new(big.Int)
	y3 = new(big.Int)
	z3 = new(big.Int)

	DoubleJacobi(
		x1, y1, z1,
		x3, y3, z3,
	)

	return x3, y3, z3
}

// Benchmarks doubling jacobi points while allocating new big.Ints for each call where Z = 1
func BenchmarkDoubleJacobi_Z1_WithAllocate(b *testing.B) {
	vector := test_vectors.JacobiDoublingVectors[0]

	for i := 0; i < b.N; i++ {
		doubleJacobiWithAllocate(
			vector.X1, vector.Y1, vector.Z1,
		)
	}
}

// Benchmarks doubling jacobi points while allocating new big.Ints for each call, where Z > 1
func BenchmarkDoubleJacobi_LargeZ_WithAllocate(b *testing.B) {
	vector := test_vectors.JacobiDoublingVectors[2]

	for i := 0; i < b.N; i++ {
		doubleJacobiWithAllocate(
			vector.X1, vector.Y1, vector.Z1,
		)
	}
}

func BenchmarkDoubleAffine(b *testing.B) {
	vector := test_vectors.JacobiDoublingVectors[0]
	x1 := new(big.Int).Set(vector.X1)
	y1 := new(big.Int).Set(vector.Y1)
	z1 := new(big.Int).Set(vector.Z1)

	ToAffine(x1, y1, z1)

	x3 := new(big.Int)
	y3 := new(big.Int)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DoubleAffine(x1, y1, x3, y3)
	}
}
