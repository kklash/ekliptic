package ekliptic

import "math/big"

// SignECDSA signs a message hash z using the private key d, and a random (or deterministically
// derived) nonce k. It returns the resulting signature parts r and s.
//
// Both the nonce k and the private key d should be generated with equal probability distribution
// over the range [1, Secp256k1_CurveOrder). SignECDSA panics if k or d is not within this range.
func SignECDSA(d, k, z *big.Int) (r, s *big.Int) {
	if !IsValidScalar(k) {
		panic("SignECDSA: expected nonce k to be in range [1, Secp256k1_CurveOrder)")
	} else if !IsValidScalar(d) {
		panic("SignECDSA: expected private key d to be in range [1, Secp256k1_CurveOrder)")
	}

	// (x, _) = k * G
	x, _ := MultiplyBasePoint(k)

	// r = x mod N
	r = new(big.Int).Mod(x, Secp256k1_CurveOrder)

	// m = rd + z
	m := x.Mul(r, d)
	m.Add(m, z)
	x = nil

	// s = k⁻¹ * m mod N
	s = InvertScalar(k)
	s.Mul(s, m)
	s.Mod(s, Secp256k1_CurveOrder)

	// always provide canonical signatures.
	//
	//  if s > (N/2):
	//    s = N - s
	if s.Cmp(Secp256k1_CurveOrderHalf) == 1 {
		s.Sub(Secp256k1_CurveOrder, s)
	}
	return
}

// VerifyECDSA returns true if the given signature (r, s) is a valid signature on message hash z
// from the given public key (pubX, pubY). Note that non-canonical ECDSA signatures (where s > N/2)
// are acceptable.
func VerifyECDSA(
	z *big.Int,
	r, s *big.Int,
	pubX, pubY *big.Int,
) bool {
	sInverse := InvertScalar(s)

	// u1 = s⁻¹ * z mod N
	u1 := new(big.Int).Mul(sInverse, z)
	u1.Mod(u1, Secp256k1_CurveOrder)

	// u2 = s⁻¹ * r mod N
	u2 := sInverse.Mul(sInverse, r)
	u2.Mod(u2, Secp256k1_CurveOrder)
	sInverse = nil

	// u1G = G * u1
	u1Gx, u1Gy := MultiplyBasePoint(u1)

	// H = (pubX, pubY)
	// u2H = H * u2
	u2Hx, u2Hy := MultiplyAffine(pubX, pubY, u2, nil)

	// P = u1G + u2H
	// px = x(P) mod N
	px, _ := AddAffine(u1Gx, u1Gy, u2Hx, u2Hy)
	px.Mod(px, Secp256k1_CurveOrder)

	return equal(r, px)
}
