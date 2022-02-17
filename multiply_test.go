package ekliptic

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestMultiplyJacobi(t *testing.T) {
	resultX := new(big.Int)
	resultY := new(big.Int)
	resultZ := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			resultX, resultY, resultZ,
			nil,
		)
		ToAffine(resultX, resultY, resultZ)

		if !equal(resultX, vector.X2) || !equal(resultY, vector.Y2) {
			t.Errorf(`jacobi multiplication failed for vector %d. Got:
	x: %x
	y: %x
Wanted:
	x: %x
	y: %x
`, i, resultX, resultY, vector.X2, vector.Y2)
		}
	}
}

func TestMultiplyJacobi_MemSafety(t *testing.T) {
	x1 := new(big.Int)
	y1 := new(big.Int)
	z1 := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		x1.Set(vector.X1)
		y1.Set(vector.Y1)
		z1.Set(one)

		MultiplyJacobi(
			x1, y1, z1,
			vector.K,
			x1, y1, z1,
			nil,
		)
		ToAffine(x1, y1, z1)

		if !equal(x1, vector.X2) || !equal(y1, vector.Y2) {
			t.Errorf(`jacobi memory-safe multiplication failed for vector %d. Got:
	x: %x
	y: %x
Wanted:
	x: %x
	y: %x
`, i, x1, y1, vector.X2, vector.Y2)
		}
	}
}

func TestMultiplyAffine(t *testing.T) {
	resultX := new(big.Int)
	resultY := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		MultiplyAffine(
			vector.X1, vector.Y1,
			vector.K,
			resultX, resultY,
			nil,
		)

		if !equal(resultX, vector.X2) || !equal(resultY, vector.Y2) {
			t.Errorf(`affine multiplication failed for vector %d. Got:
	x: %x
	y: %x
Wanted:
	x: %x
	y: %x
`, i, resultX, resultY, vector.X2, vector.Y2)
		}
	}
}

func BenchmarkMultiplyJacobi(b *testing.B) {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)

	vector := test_vectors.AffineMultiplicationVectors[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyJacobi(
			vector.X1, vector.Y1, one,
			vector.K,
			x, y, z,
			nil,
		)
	}
}

func ExampleMultiplyAffine() {
	alice, _ := new(big.Int).SetString("94a22a406a6977c1a323f23b9d7678ad08e822834d1df8adece84e30f0c25b6b", 16)
	bob, _ := new(big.Int).SetString("55ba19100104cbd2842999826e99e478efe6883ac3f3a0c7571034321e0595cf", 16)

	var alicePub, bobPub struct{ x, y big.Int }

	// derive public keys
	MultiplyBasePoint(alice, &alicePub.x, &alicePub.y)
	MultiplyBasePoint(bob, &bobPub.x, &bobPub.y)

	var yValueIsUnused big.Int

	// Alice gives Bob her public key, Bob derives the secret
	bobSharedKey := new(big.Int)
	MultiplyAffine(&alicePub.x, &alicePub.y, bob, bobSharedKey, &yValueIsUnused, nil)

	// Bob gives Alice his public key, Alice derives the secret
	aliceSharedKey := new(big.Int)
	MultiplyAffine(&bobPub.x, &bobPub.y, alice, aliceSharedKey, &yValueIsUnused, nil)

	fmt.Printf("Alice's derived secret: %x\n", aliceSharedKey)
	fmt.Printf("Bob's derived secret:   %x\n", bobSharedKey)

	// output:
	// Alice's derived secret: 375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
	// Bob's derived secret:   375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
}
