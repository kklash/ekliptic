package ekliptic

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/kklash/ekliptic/test_vectors"
)

func TestCurve(t *testing.T) {
	curve := new(Curve)

	t.Run("crypto/ecdsa", func(t *testing.T) {
		for i, vector := range test_vectors.ECDSAVectors[:20] {
			key := &ecdsa.PrivateKey{
				D: vector.PrivateKey,
				PublicKey: ecdsa.PublicKey{
					Curve: curve,
					X:     new(big.Int),
					Y:     new(big.Int),
				},
			}
			MultiplyBasePoint(key.D, key.X, key.Y)

			hash := make([]byte, 32)
			vector.Hash.FillBytes(hash)

			r, s, err := ecdsa.Sign(rand.Reader, key, hash)
			if err != nil {
				t.Errorf("failed to compute signature for vector %d: %s", i, err)
				return
			}

			if !ecdsa.Verify(&key.PublicKey, hash, r, s) {
				t.Errorf("failed to verify ECDSA signature generated using crypto/ecdsa")
				return
			}

			// Verify ekliptic's signatures also pass verification
			if !ecdsa.Verify(&key.PublicKey, hash, vector.R, vector.S) {
				t.Errorf("failed to verify ECDSA signature generated using ekliptic")
				return
			}
		}
	})

	t.Run("elliptic.Marshal", func(t *testing.T) {
		x := hexint("00cc37ea5e9e09fec6c83e5fbd7a745e3eee81d16ebd861c9e66f55518c19798")
		y := hexint("3805231b2cba4ed1b48630790489ec4b9cd44c76455856ca7c6402e0400c5d90")

		expectedPubKey, _ := hex.DecodeString(
			"04" +
				"00cc37ea5e9e09fec6c83e5fbd7a745e3eee81d16ebd861c9e66f55518c19798" +
				"3805231b2cba4ed1b48630790489ec4b9cd44c76455856ca7c6402e0400c5d90",
		)

		pubKey := elliptic.Marshal(curve, x, y)

		if !bytes.Equal(pubKey, expectedPubKey) {
			t.Errorf("Unexpected uncompressed public key:\nWanted '%x'\nGot    '%x'", expectedPubKey, pubKey)
			return
		}

		parsedX, parsedY := elliptic.Unmarshal(curve, pubKey)

		if !equal(parsedX, x) || !equal(parsedY, y) {
			t.Errorf(`unmarshaled unexpected public key coordinates:
  x: %.64x
  y: %.64x
Wanted:
  x: %.64x
  y: %.64x
`, parsedX, parsedY, x, y)
			return
		}
	})

	t.Run("IsOnCurve", func(t *testing.T) {
		if curve.IsOnCurve(zero, zero) {
			t.Errorf("expected infinity not to be considered on curve")
			return
		}
	})

	t.Run("Params", func(*testing.T) {
		p := curve.Params().P
		p.Sub(p, one)
		if equal(Secp256k1_P, p) {
			t.Errorf("expected CurveParams to be independent of ekliptic constants")
			return
		}
	})
}
