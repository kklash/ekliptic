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
