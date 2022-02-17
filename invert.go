package ekliptic

import (
	"fmt"
	"math/big"
)

func invert(n *big.Int) {
	i := n.ModInverse(n, Secp256k1_P)
	if i == nil {
		err := fmt.Sprintf("cannot take multiplicative inverse of value: %s - P is probably not prime", n.Text(16))
		panic(err)
	}
}
