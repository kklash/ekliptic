package ekliptic

import "math/big"

func mod(n *big.Int) {
	n.Mod(n, Secp256k1_P)
}
