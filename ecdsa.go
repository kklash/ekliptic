package ekliptic

import "math/big"

// SignECDSA signs a message hash z using the private key d, and a random (or deterministically
// derived) nonce k. It sets r and s to the resulting signature parts.
//
// Both the nonce k and the private key d should be generated with equal probability distribution
// over the range [1, Secp256k1_CurveOrder). SignECDSA Panics if k or d is not within this range.
func SignECDSA(
	d, k, z *big.Int,
	r, s *big.Int,
) {
	if !IsValidScalar(k) {
		panic("SignECDSA: expected nonce k to be in range [1, Secp256k1_CurveOrder)")
	} else if !IsValidScalar(d) {
		panic("SignECDSA: expected private key d to be in range [1, Secp256k1_CurveOrder)")
	}

	// (x, _) = k * G
	x := new(big.Int)
	MultiplyBasePoint(k, x, new(big.Int))

	// r = x mod N
	r.Mod(x, Secp256k1_CurveOrder)

	// m = rd + z
	m := x.Mul(r, d)
	m.Add(m, z)
	x = nil

	// s = k⁻¹ * m mod N
	s.ModInverse(k, Secp256k1_CurveOrder)
	s.Mul(s, m)
	s.Mod(s, Secp256k1_CurveOrder)

	// always provide canonical signatures.
	//
	//  if s > (N/2):
	//    s = N - s
	if s.Cmp(Secp256k1_CurveOrderHalf) == 1 {
		s.Sub(Secp256k1_CurveOrder, s)
	}
}

// VerifyECDSA returns true if the given signature (r, s) is a valid signature on message hash z
// from the given public key (pubX, pubY). Note that non-canonical ECDSA signatures (where s > N/2)
// are acceptable.
func VerifyECDSA(
	z *big.Int,
	r, s *big.Int,
	pubX, pubY *big.Int,
) bool {
	sInverse := new(big.Int).ModInverse(s, Secp256k1_CurveOrder)

	// u1 = s⁻¹ * z mod N
	u1 := new(big.Int).Mul(sInverse, z)
	u1.Mod(u1, Secp256k1_CurveOrder)

	// u2 = s⁻¹ * r mod N
	u2 := sInverse.Mul(sInverse, r)
	u2.Mod(u2, Secp256k1_CurveOrder)
	sInverse = nil

	// u1G = G * u1
	u1Gx := u1
	u1Gy := new(big.Int)
	MultiplyBasePoint(u1, u1Gx, u1Gy)
	u1 = nil

	// H = (pubX, pubY)
	// u2H = H * u2
	u2Hx := u2
	u2Hy := new(big.Int)
	MultiplyAffine(pubX, pubY, u2, u2Hx, u2Hy, nil)
	u2 = nil

	// P = u1G + u2H
	// px = x(p) mod N
	px := u1Gx
	AddAffine(u1Gx, u1Gy, u2Hx, u2Hy, px, u1Gy)
	px.Mod(px, Secp256k1_CurveOrder)
	u1Gx = nil
	u1Gy = nil

	return equal(r, px)
}
