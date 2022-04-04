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
	if k.Cmp(Secp256k1_CurveOrder) >= 0 || k.Cmp(one) == -1 {
		panic("SignECDSA: expected nonce k to be in range [1, Secp256k1_CurveOrder)")
	} else if d.Cmp(Secp256k1_CurveOrder) >= 0 || d.Cmp(one) == -1 {
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
