package test_vectors

import (
	_ "embed"
	"encoding/json"
	"math/big"
)

// JacobiPointVector represents a Jacobian point and the affine coordinates it equates to.
// One affine point can be expressed in many different Jacobian points.
type JacobiPointVector struct {
	JacobiX, JacobiY, JacobiZ *big.Int
	X, Y                      *big.Int
}

//go:embed jacobi_point.json
var jacobiPointJsonBytes []byte

func loadJacobiPointVectors() ([]*JacobiPointVector, error) {
	var rawJsonObjects []map[string]string

	if err := json.Unmarshal(jacobiPointJsonBytes, &rawJsonObjects); err != nil {
		return nil, err
	}

	vectors := make([]*JacobiPointVector, len(rawJsonObjects))

	for i, obj := range rawJsonObjects {
		vectors[i] = &JacobiPointVector{
			JacobiX: hexint(obj["jacobiX"]),
			JacobiY: hexint(obj["jacobiY"]),
			JacobiZ: hexint(obj["jacobiZ"]),
			X:       hexint(obj["x"]),
			Y:       hexint(obj["y"]),
		}
	}

	return vectors, nil
}
