package ekliptic

import (
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestNegate(t *testing.T) {
	for i, vector := range test_vectors.NegatedPointVectors {
		y := Negate(vector.EvenY)

		if !equal(y, vector.OddY) {
			t.Errorf("negation from even to odd failed for vector %d.\nWanted 0x%.64x\n    Got 0x%.64x", i, vector.OddY, y)
		}

		y = Negate(y)

		if !equal(y, vector.EvenY) {
			t.Errorf("negation from odd to even failed for vector %d.\nWanted 0x%.64x\n    Got 0x%.64x", i, vector.EvenY, y)
		}
	}

	// Negating zero should be zero, because 0 - 0 = 0
	if !equal(Negate(zero), zero) {
		t.Errorf("expected negating zero to result in zero")
	}
}

func BenchmarkNegate(b *testing.B) {
	vector := test_vectors.NegatedPointVectors[0]
	for i := 0; i < b.N; i++ {
		Negate(vector.EvenY)
	}
}
