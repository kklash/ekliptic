package ekliptic

import (
	"crypto/rand"
	"io"
	"math/big"
)

// NewPrivateKey generates a random private key from the given source of randomness.
// Ensures the distribution of possible private keys is uniformly distributed from [0..N-1].
func NewPrivateKey(random io.Reader) (*big.Int, error) {
	return rand.Int(random, Secp256k1_CurveOrder)
}
