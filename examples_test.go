package ekliptic_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	mathrand "math/rand"

	"github.com/kklash/ekliptic"
)

// *** ATTENTION ***
// *****************
// Modifying this file? Make sure to copy the changes to the README's examples section!

// Generate a public key from a private key.
func ExampleMultiplyBasePoint() {
	privateKey, _ := new(big.Int).SetString("c370af8c091812ef7f6bfaffb494b1046fb25486c9873243b80826daef3ec583", 16)
	x := new(big.Int)
	y := new(big.Int)

	ekliptic.MultiplyBasePoint(privateKey, x, y)

	fmt.Println("Public key:")
	fmt.Printf(" x: %x\n", x)
	fmt.Printf(" y: %x\n", y)

	// output:
	// Public key:
	//  x: 76cd66c6cca75278ff408ce67290537367719154ae2b96448327fe4033ddcfc7
	//  y: 35663ecbb64397bb9bd79155a1e6b138c2fb8fa1f11355f8e9e97ddd88a78e49
}

// Construct an ECDH shared secret.
func ExampleMultiplyAffine() {
	alicePriv, _ := new(big.Int).SetString("94a22a406a6977c1a323f23b9d7678ad08e822834d1df8adece84e30f0c25b6b", 16)
	bobPriv, _ := new(big.Int).SetString("55ba19100104cbd2842999826e99e478efe6883ac3f3a0c7571034321e0595cf", 16)

	var alicePub, bobPub struct{ x, y big.Int }

	// derive public keys
	ekliptic.MultiplyBasePoint(alicePriv, &alicePub.x, &alicePub.y)
	ekliptic.MultiplyBasePoint(bobPriv, &bobPub.x, &bobPub.y)

	var yValueIsUnused big.Int

	// Alice gives Bob her public key, Bob derives the secret
	bobSharedKey := new(big.Int)
	ekliptic.MultiplyAffine(&alicePub.x, &alicePub.y, bobPriv, bobSharedKey, &yValueIsUnused, nil)

	// Bob gives Alice his public key, Alice derives the secret
	aliceSharedKey := new(big.Int)
	ekliptic.MultiplyAffine(&bobPub.x, &bobPub.y, alicePriv, aliceSharedKey, &yValueIsUnused, nil)

	fmt.Printf("Alice's derived secret: %x\n", aliceSharedKey)
	fmt.Printf("Bob's derived secret:   %x\n", bobSharedKey)

	// output:
	// Alice's derived secret: 375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
	// Bob's derived secret:   375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
}

// Sign a message digest.
func ExampleSignECDSA() {
	randReader := mathrand.New(mathrand.NewSource(1))

	key, _ := ekliptic.NewPrivateKey(randReader)

	// This could also come from RFC6979 (github.com/kklash/rfc6979)
	nonce, _ := cryptorand.Int(randReader, ekliptic.Secp256k1_CurveOrder)

	hashedMessage := sha256.Sum256([]byte("i love you"))
	hashedMessageInt := new(big.Int).SetBytes(hashedMessage[:])

	r := new(big.Int)
	s := new(big.Int)

	ekliptic.SignECDSA(
		key, nonce, hashedMessageInt,
		r, s,
	)

	fmt.Printf("r: %x\n", r)
	fmt.Printf("s: %x\n", s)

	var pub struct{ x, y big.Int }
	ekliptic.MultiplyBasePoint(key, &pub.x, &pub.y)

	valid := ekliptic.VerifyECDSA(hashedMessageInt, r, s, &pub.x, &pub.y)
	fmt.Printf("valid: %v\n", valid)

	// output:
	//
	// r: 4a821d5ec008712983929de448b8afb6c24e5a1b97367b9a65b6220d7f083fe3
	// s: 381d053be61243d950865d7b8eb6b5ba48fbabfe7fda81af3183a184a02f5d51
	// valid: true
}

// Find possible Y-coordinates for an X. Used to uncompress a public key, where
// you may only have the full X-coordinate of the public key.
func ExampleWeierstrass() {
	compressedKey, _ := hex.DecodeString("030000000000000000000000000000000000000000000000000000000000000001")

	publicKeyX := new(big.Int).SetBytes(compressedKey[1:])
	evenY, oddY := ekliptic.Weierstrass(publicKeyX)

	var publicKeyY *big.Int
	if compressedKey[0]%2 == 0 {
		publicKeyY = evenY
	} else {
		publicKeyY = oddY
	}

	fmt.Println("uncompressed key:")
	fmt.Printf("x: %.64x\n", publicKeyX)
	fmt.Printf("y: %.64x\n", publicKeyY)

	// output:
	// uncompressed key:
	// x: 0000000000000000000000000000000000000000000000000000000000000001
	// y: bde70df51939b94c9c24979fa7dd04ebd9b3572da7802290438af2a681895441
}

func ExampleCurve() {
	d, _ := new(big.Int).SetString("18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725", 16)
	key := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: new(ekliptic.Curve),
			X:     new(big.Int),
			Y:     new(big.Int),
		},
	}

	// Compute the public key
	ekliptic.MultiplyBasePoint(key.D, key.X, key.Y)

	hashedMessage := sha256.Sum256([]byte("i love you"))

	r, s, err := ecdsa.Sign(rand.Reader, key, hashedMessage[:])
	if err != nil {
		panic("failed to compute signature: " + err.Error())
	}

	if ecdsa.Verify(&key.PublicKey, hashedMessage[:], r, s) {
		fmt.Println("verified ECDSA signature using crypto/ecdsa")
	}

	// output:
	// verified ECDSA signature using crypto/ecdsa
}
