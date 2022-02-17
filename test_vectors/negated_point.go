package test_vectors

import (
	_ "embed"
	"encoding/json"
	"math/big"
)

// NegatedPointVector represents an affine X-coordinate and the two possible
// Y-coordinates on the curve at that X-coordinate: one even and one odd.
type NegatedPointVector struct {
	X, EvenY, OddY *big.Int
}

//go:embed negated_point.json
var negatedPointJsonBytes []byte

func loadNegatedPointVectors() ([]*NegatedPointVector, error) {
	var rawJsonObjects []map[string]string

	if err := json.Unmarshal(negatedPointJsonBytes, &rawJsonObjects); err != nil {
		return nil, err
	}

	vectors := make([]*NegatedPointVector, len(rawJsonObjects))

	for i, obj := range rawJsonObjects {
		vectors[i] = &NegatedPointVector{
			X:     hexint(obj["x"]),
			EvenY: hexint(obj["evenY"]),
			OddY:  hexint(obj["oddY"]),
		}
	}

	return vectors, nil
}
