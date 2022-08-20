package ekliptic

import (
	"crypto/rand"
	"io"
	"math/big"
)

var curveOrderMinusOne = new(big.Int).Sub(Secp256k1_CurveOrder, one)

// RandomScalar securely generates a random scalar value from the given source of randomness.
// Ensures the distribution of possible scalars is uniformly distributed from [1..N-1].
//
// Use this function to generate private keys.
func RandomScalar(random io.Reader) (*big.Int, error) {
	r, err := rand.Int(random, curveOrderMinusOne)
	if err != nil {
		return nil, err
	}
	return r.Add(r, one), nil
}

// IsValidScalar returns true if the given integer is a valid secp256k1 scalar (private
// key), i.e. a number in the range [1, N-1] where N is the secp256k1 curve order.
func IsValidScalar(d *big.Int) bool {
	return d.Cmp(zero) == 1 && d.Cmp(Secp256k1_CurveOrder) == -1
}
