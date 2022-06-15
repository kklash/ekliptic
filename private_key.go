package ekliptic

import (
	"crypto/rand"
	"io"
	"math/big"
)

// NewPrivateKey generates a random private key from the given source of randomness.
// Ensures the distribution of possible private keys is uniformly distributed from [1..N-1].
func NewPrivateKey(random io.Reader) (*big.Int, error) {
	r, err := rand.Int(random, new(big.Int).Sub(Secp256k1_CurveOrder, one))
	if err != nil {
		return nil, err
	}
	return r.Add(r, one), nil
}
