package ekliptic

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestWeierstrass(t *testing.T) {
	for i, vector := range test_vectors.NegatedPointVectors {
		evenY, oddY := Weierstrass(vector.X)

		if !equal(evenY, vector.EvenY) || !equal(oddY, vector.OddY) {
			t.Errorf(`Weierstrass calculation failed for vector %d. Wanted
	evenY: %x
	oddY:  %x
Got:
	evenY: %x
	oddY:  %x
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
		func() {
			defer func() {
				panicValue := recover()
				if panicValue == nil {
					t.Errorf("expected panic when calling Weierstrass on invalid x value: '%x'", x)
				}
			}()

			Weierstrass(x)
		}()
	}
}

func ExampleWeierstrass() {
	evenY, oddY := Weierstrass(big.NewInt(1))

	fmt.Printf("even: %x\n", evenY)
	fmt.Printf("odd:  %x\n", oddY)

	// output:
	// even: 4218f20ae6c646b363db68605822fb14264ca8d2587fdd6fbc750d587e76a7ee
	// odd:  bde70df51939b94c9c24979fa7dd04ebd9b3572da7802290438af2a681895441
}
