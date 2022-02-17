package test_vectors

import (
	_ "embed"
	"encoding/json"
	"math/big"
)

// JacobiAdditionVector represents the result of a addition operation using
// Jacobian coordinates, where P1 + P2 = P3.
type JacobiAdditionVector struct {
	X1, Y1, Z1 *big.Int
	X2, Y2, Z2 *big.Int
	X3, Y3, Z3 *big.Int
}

//go:embed jacobi_addition.json
var jacobiAdditionJsonBytes []byte

func loadJacobiAdditionVectors() ([]*JacobiAdditionVector, error) {
	var rawJsonObjects []map[string]string

	if err := json.Unmarshal(jacobiAdditionJsonBytes, &rawJsonObjects); err != nil {
		return nil, err
	}

	vectors := make([]*JacobiAdditionVector, len(rawJsonObjects))

	for i, obj := range rawJsonObjects {
		vectors[i] = &JacobiAdditionVector{
			X1: hexint(obj["x1"]),
			Y1: hexint(obj["y1"]),
			Z1: hexint(obj["z1"]),

			X2: hexint(obj["x2"]),
			Y2: hexint(obj["y2"]),
			Z2: hexint(obj["z2"]),

			X3: hexint(obj["x3"]),
			Y3: hexint(obj["y3"]),
			Z3: hexint(obj["z3"]),
		}
	}

	return vectors, nil
}
