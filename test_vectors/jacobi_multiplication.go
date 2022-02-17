package test_vectors

import (
	_ "embed"
	"encoding/json"
	"math/big"
)

// AffineMultiplicationVector represents the result of multiplying an affine point by some scalar value k.
type AffineMultiplicationVector struct {
	X1, Y1 *big.Int
	K      *big.Int
	X2, Y2 *big.Int
}

//go:embed jacobi_multiplication.json
var affineMultiplicationJsonBytes []byte

func loadAffineMultiplicationVectors() ([]*AffineMultiplicationVector, error) {
	var rawJsonObjects []map[string]string

	if err := json.Unmarshal(affineMultiplicationJsonBytes, &rawJsonObjects); err != nil {
		return nil, err
	}

	vectors := make([]*AffineMultiplicationVector, len(rawJsonObjects))

	for i, obj := range rawJsonObjects {
		vectors[i] = &AffineMultiplicationVector{
			X1: hexint(obj["x1"]),
			Y1: hexint(obj["y1"]),
			K:  hexint(obj["k"]),
			X2: hexint(obj["x2"]),
			Y2: hexint(obj["y2"]),
		}
	}

	return vectors, nil
}
