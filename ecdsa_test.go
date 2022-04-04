package ekliptic

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestSignECDSA(t *testing.T) {
	r := new(big.Int)
	s := new(big.Int)

	for i, vector := range test_vectors.ECDSAVectors {
		SignECDSA(
			vector.PrivateKey, vector.Nonce, vector.Hash,
			r, s,
		)

		if !equal(r, vector.R) || !equal(s, vector.S) {
			t.Errorf(`invalid ECDSA signature for vector %d. Got:
	r: %.64x
	s: %.64x
Wanted:
	r: %.64x
	s: %.64x
`, i, r, s, vector.R, vector.S)
		}
	}
}

func BenchmarkSignECDSA(b *testing.B) {
	r := new(big.Int)
	s := new(big.Int)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector := test_vectors.ECDSAVectors[i%len(test_vectors.ECDSAVectors)]
		SignECDSA(
			vector.PrivateKey, vector.Nonce, vector.Hash,
			r, s,
		)
	}
}

func ExampleSignECDSA() {
	randReader := mathrand.New(mathrand.NewSource(1))

	key, _ := cryptorand.Int(randReader, Secp256k1_CurveOrder)
	nonce, _ := cryptorand.Int(randReader, Secp256k1_CurveOrder)

	hashedMessage := sha256.Sum256([]byte("i love you"))
	hashedMessageInt := new(big.Int).SetBytes(hashedMessage[:])

	r := new(big.Int)
	s := new(big.Int)

	SignECDSA(
		key, nonce, hashedMessageInt,
		r, s,
	)

	fmt.Printf("r: %x\n", r)
	fmt.Printf("s: %x\n", s)

	// output:
	//
	// r: 4a821d5ec008712983929de448b8afb6c24e5a1b97367b9a65b6220d7f083fe3
	// s: 2e4f380e0ea1dfcb7cced430437c98b4570a06b3e929a3b19e6bbd53df2cf3f6
}
