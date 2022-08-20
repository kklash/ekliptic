package ekliptic

import "math/big"

func modCoordinate(n *big.Int) {
	n.Mod(n, Secp256k1_P)
}
