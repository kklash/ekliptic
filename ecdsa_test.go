package ekliptic

import (
	"math/big"
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
