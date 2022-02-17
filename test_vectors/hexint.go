package test_vectors

import (
	"fmt"
	"math/big"
)

func hexint(s string) *big.Int {
	i, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(fmt.Sprintf("Failed to parse hexint: '%s'", s))
	}
	return i
}
