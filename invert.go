package ekliptic

import (
	"fmt"
	"math/big"
)

// invertCoordinate inverts the given point coordinate value n modulo the
// curve's prime parameter.
func invertCoordinate(n *big.Int) {
	i := n.ModInverse(n, Secp256k1_P)
	if i == nil {
		err := fmt.Sprintf("cannot take multiplicative inverse of value: %s - P is probably not prime", n.Text(16))
		panic(err)
	}
}

// InvertScalar returns d⁻¹: the multiplicative inverse of the given scalar value d,
// modulo the curve order N:
//  d * d⁻¹ = 1 mod N
//
// Multiplying any point A by both d and d⁻¹ will return the same point A:
//  d * A * d⁻¹ = A
func InvertScalar(d *big.Int) *big.Int {
	return new(big.Int).ModInverse(d, Secp256k1_CurveOrder)
}
