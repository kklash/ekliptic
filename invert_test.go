package ekliptic

import (
	"crypto/rand"
	"testing"
)

func TestInvertScalar(t *testing.T) {
	for i := 0; i < 10; i++ {
		r, err := RandomScalar(rand.Reader)
		if err != nil {
			t.Errorf("Failed to generate random scalar: %s", err)
			return
		}

		x, y := MultiplyBasePoint(r)
		x, y = MultiplyAffine(x, y, InvertScalar(r), nil)

		if !EqualAffine(x, y, Secp256k1_GeneratorX, Secp256k1_GeneratorY) {
			t.Errorf("expected to get generator point back after multiplying public key by inverse private key")
			return
		}
	}
}

func BenchmarkInvertScalar(b *testing.B) {
	r, err := RandomScalar(rand.Reader)
	if err != nil {
		b.Errorf("Failed to generate random scalar: %s", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InvertScalar(r)
	}
}
