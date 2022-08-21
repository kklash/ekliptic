package ekliptic

import (
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestSignECDSA(t *testing.T) {
	for i, vector := range test_vectors.ECDSAVectors {
		r, s := SignECDSA(vector.PrivateKey, vector.Nonce, vector.Hash)

		if !equal(r, vector.R) || !equal(s, vector.S) {
			t.Errorf(`invalid ECDSA signature for vector %d. Got:
	r: %.64x
	s: %.64x
Wanted:
	r: %.64x
	s: %.64x
`, i, r, s, vector.R, vector.S)
			return
		}

		pubX, pubY := MultiplyBasePoint(vector.PrivateKey)

		if !VerifyECDSA(vector.Hash, vector.R, vector.S, pubX, pubY) {
			t.Errorf("failed to verify ECDSA signature for vector %d", i)
			return
		}
	}
}

func BenchmarkSignECDSA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vector := test_vectors.ECDSAVectors[i%len(test_vectors.ECDSAVectors)]
		SignECDSA(vector.PrivateKey, vector.Nonce, vector.Hash)
	}
}

func BenchmarkVerifyECDSA(b *testing.B) {
	vector := test_vectors.ECDSAVectors[4]
	pubX, pubY := MultiplyBasePoint(vector.PrivateKey)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		VerifyECDSA(vector.Hash, vector.R, vector.S, pubX, pubY)
	}
}
