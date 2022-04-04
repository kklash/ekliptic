package test_vectors

import (
	_ "embed"
	"encoding/json"
	"math/big"
)

// ECDSAVector represents a test vector for the elliptic curve digital signature algorithm
// on a given message hash with a private key and nonce.
type ECDSAVector struct {
	Hash       *big.Int
	PrivateKey *big.Int
	Nonce      *big.Int
	R, S       *big.Int
}

//go:embed ecdsa.json
var ecdsaJsonBytes []byte

func loadECDSAVectors() ([]*ECDSAVector, error) {
	var rawJsonObjects []map[string]string

	if err := json.Unmarshal(ecdsaJsonBytes, &rawJsonObjects); err != nil {
		return nil, err
	}

	vectors := make([]*ECDSAVector, len(rawJsonObjects))

	for i, obj := range rawJsonObjects {
		vectors[i] = &ECDSAVector{
			Hash:       hexint(obj["hash"]),
			PrivateKey: hexint(obj["d"]),
			Nonce:      hexint(obj["nonce"]),
			R:          hexint(obj["r"]),
			S:          hexint(obj["s"]),
		}
	}

	return vectors, nil
}
