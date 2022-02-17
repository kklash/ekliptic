package test_vectors

import (
	_ "embed"
	"encoding/json"
	"math/big"
)

// JacobiDoublingVector represents the result of doubling a Jacobian point, where P1 * 2 = P3
type JacobiDoublingVector struct {
	X1, Y1, Z1 *big.Int
	X3, Y3, Z3 *big.Int
}

//go:embed jacobi_doubling.json
var jacobiDoublingJsonBytes []byte

func loadJacobiDoublingVectors() ([]*JacobiDoublingVector, error) {
	var rawJsonObjects []map[string]string

	if err := json.Unmarshal(jacobiDoublingJsonBytes, &rawJsonObjects); err != nil {
		return nil, err
	}

	vectors := make([]*JacobiDoublingVector, len(rawJsonObjects))

	for i, obj := range rawJsonObjects {
		vectors[i] = &JacobiDoublingVector{
			X1: hexint(obj["x1"]),
			Y1: hexint(obj["y1"]),
			Z1: hexint(obj["z1"]),

			X3: hexint(obj["x3"]),
			Y3: hexint(obj["y3"]),
			Z3: hexint(obj["z3"]),
		}
	}

	return vectors, nil
}
