package ekliptic

import (
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestWeierstrass(t *testing.T) {
	for i, vector := range test_vectors.NegatedPointVectors {
		evenY, oddY := Weierstrass(vector.X)

		if !equal(evenY, vector.EvenY) || !equal(oddY, vector.OddY) {
			t.Errorf(`Weierstrass calculation failed for vector %d. Wanted
	evenY: %.64x
	oddY:  %.64x
Got:
	evenY: %.64x
	oddY:  %.64x
`, i, vector.EvenY, vector.OddY, evenY, oddY)
		}
	}
}

func TestWeierstrass_NotOnCurve(t *testing.T) {
	invalidXs := []*big.Int{
		hexint("b906eae4400782607f482c77b0c3c8e049577d8c1ff0779374818a3a2f5a3a34"),
		hexint("4323d9bc9c1c255f256e6549828aa2e40052325cd5eb277f5836d8a5713aac1d"),
		big.NewInt(0),
		Secp256k1_P,
		new(big.Int).Add(Secp256k1_P, one),
	}

	for _, x := range invalidXs {
		evenY, oddY := Weierstrass(x)
		if evenY != nil || oddY != nil {
			t.Errorf("expected Weierstrass on invalid x value to return nil: '%.64x'", x)
			return
		}
	}
}

func BenchmarkWeierstrass(b *testing.B) {
	vector := test_vectors.NegatedPointVectors[0]

	for i := 0; i < b.N; i++ {
		Weierstrass(vector.X)
	}
}
