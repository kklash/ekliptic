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
	x, y := ekliptic.MultiplyBasePoint(privateKey)

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

	var alicePub, bobPub struct{ x, y *big.Int }

	// derive public keys
	alicePub.x, alicePub.y = ekliptic.MultiplyBasePoint(alicePriv)
	bobPub.x, bobPub.y = ekliptic.MultiplyBasePoint(bobPriv)

	// Alice gives Bob her public key, Bob derives the secret
	bobSharedKey, _ := ekliptic.MultiplyAffine(alicePub.x, alicePub.y, bobPriv, nil)

	// Bob gives Alice his public key, Alice derives the secret
	aliceSharedKey, _ := ekliptic.MultiplyAffine(bobPub.x, bobPub.y, alicePriv, nil)

	fmt.Printf("Alice's derived secret: %x\n", aliceSharedKey)
	fmt.Printf("Bob's derived secret:   %x\n", bobSharedKey)

	// output:
	// Alice's derived secret: 375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
	// Bob's derived secret:   375a5d26649704863562930ded2193a0569f90f4eb4e63f0fee72c4c05268feb
}

// Sign a message digest.
func ExampleSignECDSA() {
	randReader := mathrand.New(mathrand.NewSource(1))

	key, _ := ekliptic.RandomScalar(randReader)

	// This could also come from RFC6979 (github.com/kklash/rfc6979)
	nonce, _ := cryptorand.Int(randReader, ekliptic.Secp256k1_CurveOrder)

	hashedMessage := sha256.Sum256([]byte("i love you"))
	hashedMessageInt := new(big.Int).SetBytes(hashedMessage[:])

	r, s := ekliptic.SignECDSA(key, nonce, hashedMessageInt)

	fmt.Printf("r: %x\n", r)
	fmt.Printf("s: %x\n", s)

	var pub struct{ x, y *big.Int }
	pub.x, pub.y = ekliptic.MultiplyBasePoint(key)

	valid := ekliptic.VerifyECDSA(hashedMessageInt, r, s, pub.x, pub.y)
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
	pubX, pubY := ekliptic.MultiplyBasePoint(d)
	key := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: new(ekliptic.Curve),
			X:     pubX,
			Y:     pubY,
		},
	}

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

// InvertScalar is useful for reversibly blinding a value you don't want to reveal.
// Alice can blind any point A with some random scalar s to produce a blinded point B:
//  B = s * A
//
// B can then be blinded by another party Bob, with their own secret r:
//  C = r * B
//
// Alice can then unblind C by inverting their secret s:
//  M = s⁻¹ * C
//  M = s⁻¹ * (r * s * A)
//  M = r * A
//
// Alice now knows r * A without knowing r, having revealed neither A nor the final result M to Bob.
func ExampleInvertScalar() {
	// Alice: input parameters
	s, _ := new(big.Int).SetString("2fc9374cad648e33f78dd294578dd960281e05744b27faa1ffe1e7175bd6901d", 16)
	aX, _ := new(big.Int).SetString("8bf20851cc16007dbf3df0c109dc016b360ca0f729f368ea38c385ceeffaf3cf", 16)
	aY, _ := new(big.Int).SetString("a0c7bd73154c02cc5e002c3a4f876158a4276c185ef859df589675a92c745e3a", 16)

	// Bob: input parameters
	r, _ := new(big.Int).SetString("825e7984ae7843f9c13371d9a54143a465b1e2d278e67de1ca713127e40a52f1", 16)

	// Alice: blind the input with secret s send the result B to Bob.
	// B = s * A
	bX, bY := ekliptic.MultiplyAffine(aX, aY, s, nil)

	// Bob: receive A, blind it again with secret r and return result C to Alice.
	// C = r * B = r * s * G
	cX, cY := ekliptic.MultiplyAffine(bX, bY, r, nil)

	// Alice: unblind C with the inverse of s:
	// s⁻¹ * C = s⁻¹ * r * s * G = r * G
	sInv := ekliptic.InvertScalar(s)

	mX, mY := ekliptic.MultiplyAffine(cX, cY, sInv, nil)

	expectedX, expectedY := ekliptic.MultiplyAffine(aX, aY, r, nil)

	if !ekliptic.EqualAffine(mX, mY, expectedX, expectedY) {
		fmt.Println("did not find expected unblinded point")
		return
	}

	fmt.Println("found correct unblinded point q * h * G")

	// output:
	// found correct unblinded point q * h * G
}
