package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestIsOnCurveJacobi(t *testing.T) {
	for i, vector := range test_vectors.JacobiPointVectors {
		x := new(big.Int).Set(vector.JacobiX)
		y := new(big.Int).Set(vector.JacobiY)
		z := new(big.Int).Set(vector.JacobiZ)

		if !IsOnCurveJacobi(x, y, z) {
			t.Errorf("on-curve check failed for valid jacobi point vector #%d", i)
		}

		switch i % 3 {
		case 0:
			x.Add(x, one)
		case 1:
			y.Mul(y, two)
			modCoordinate(y)
		case 2:
			z.Sub(z, one)
		}

		if IsOnCurveJacobi(x, y, z) {
			t.Errorf("on-curve check failed for INVALID jacobi point vector #%d", i)
		}
	}
}

func TestIsOnCurveAffine(t *testing.T) {
	for i, vector := range test_vectors.JacobiPointVectors {
		x := new(big.Int).Set(vector.JacobiX)
		y := new(big.Int).Set(vector.JacobiY)
		z := new(big.Int).Set(vector.JacobiZ)

		ToAffine(x, y, z)

		if !IsOnCurveAffine(x, y) {
			t.Errorf("on-curve check failed for valid affine point vector #%d", i)
		}

		switch i % 2 {
		case 0:
			x.Add(x, one)
		case 1:
			y.Mul(y, two)
			modCoordinate(y)
		}

		if IsOnCurveAffine(x, y) {
			t.Errorf("on-curve check failed for INVALID affine point vector #%d", i)
		}
	}
}

func BenchmarkIsOnCurveJacobi(b *testing.B) {
	vector := test_vectors.JacobiPointVectors[0]

	for i := 0; i < b.N; i++ {
		IsOnCurveJacobi(vector.JacobiX, vector.JacobiY, vector.JacobiZ)
	}
}

func BenchmarkIsOnCurveAffine(b *testing.B) {
	vector := test_vectors.JacobiPointVectors[0]
	x := new(big.Int).Set(vector.JacobiX)
	y := new(big.Int).Set(vector.JacobiY)
	z := new(big.Int).Set(vector.JacobiZ)
	ToAffine(x, y, z)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsOnCurveAffine(x, y)
	}
}
