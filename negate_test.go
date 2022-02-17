package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestNegate(t *testing.T) {
	y := new(big.Int)
	for i, vector := range test_vectors.NegatedPointVectors {
		y.Set(vector.EvenY)

		Negate(y)

		if !equal(y, vector.OddY) {
			t.Errorf("negation from even to odd failed for vector %d.\nWanted 0x%x\n    Got 0x%x", i, vector.OddY, y)
		}

		Negate(y)

		if !equal(y, vector.EvenY) {
			t.Errorf("negation from odd to even failed for vector %d.\nWanted 0x%x\n    Got 0x%x", i, vector.EvenY, y)
		}
	}

	// Negating zero should be zero, because 0 - 0 = 0
	y.Set(zero)
	Negate(y)
	if !equal(y, zero) {
		t.Errorf("expected negating zero to result in zero")
	}
}

func BenchmarkNegate(b *testing.B) {
	vector := test_vectors.NegatedPointVectors[0]
	y := new(big.Int).Set(vector.EvenY)

	for i := 0; i < b.N; i++ {
		Negate(y)
	}
}
