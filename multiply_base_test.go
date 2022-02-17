package ekliptic

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestMultiplyBasePoint(t *testing.T) {
	x := new(big.Int)
	y := new(big.Int)

	for i, vector := range test_vectors.AffineMultiplicationVectors {
		if !equal(vector.X1, Secp256k1_GeneratorX) ||
			!equal(vector.Y1, Secp256k1_GeneratorY) {
			continue
		}

		MultiplyBasePoint(vector.K, x, y)

		if !equal(x, vector.X2) || !equal(y, vector.Y2) {
			t.Errorf(`multiplying base point failed for vector %d. Got
	x: %x
	y: %x
Wanted:
	x: %x
	y: %x
	`, i, x, y, vector.X2, vector.Y2)
		}
	}
}

func BenchmarkMultiplyBasePoint(b *testing.B) {
	vector := test_vectors.AffineMultiplicationVectors[0]
	x := new(big.Int)
	y := new(big.Int)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiplyBasePoint(vector.K, x, y)
	}
}

func ExampleMultiplyBasePoint() {
	// Generate a public key from a private key.
	privateKey, _ := new(big.Int).SetString("c370af8c091812ef7f6bfaffb494b1046fb25486c9873243b80826daef3ec583", 16)
	x := new(big.Int)
	y := new(big.Int)

	MultiplyBasePoint(privateKey, x, y)

	fmt.Println("Public key:")
	fmt.Printf(" x: %x\n", x)
	fmt.Printf(" y: %x\n", y)

	// output:
	// Public key:
	//  x: 76cd66c6cca75278ff408ce67290537367719154ae2b96448327fe4033ddcfc7
	//  y: 35663ecbb64397bb9bd79155a1e6b138c2fb8fa1f11355f8e9e97ddd88a78e49
}
