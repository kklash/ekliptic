package ekliptic

import (
	"math/big"
)

// AddJacobi adds two Jacobian coordinate points on the secp256k1 curve:
//
//  P1 + P2 = P3
//  (x1, y1, z1) + (x2, y2, z2) = (x3, y3, z3)
//
// It stores the resulting Jacobian point in x3, y3, and z3.
//
// We use the "add-1998-cmo-2" addition formulas.
//
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-1998-cmo-2
//
//  Z1Z1 = Z1²
//  Z2Z2 = Z2²
//  U1 = X1*Z2Z2
//  U2 = X2*Z1Z1
//  S1 = Y1*Z2*Z2Z2
//  S2 = Y2*Z1*Z1Z1
//  H = U2-U1
//  HH = H²
//  HHH = H*HH
//  r = S2-S1
//  V = U1*HH
//  X3 = r²-HHH-2*V
//  Y3 = r*(V-X3)-S1*HHH
//  Z3 = Z1*Z2*H
func AddJacobi(
	x1, y1, z1 *big.Int, // P1
	x2, y2, z2 *big.Int, // P2
	x3, y3, z3 *big.Int, // P3
) {
	if equal(x1, zero) || equal(y1, zero) {
		// P1 == 0: return P2
		x3.Set(x2)
		y3.Set(y2)
		z3.Set(z2)
		return
	}
	if equal(x2, zero) || equal(y2, zero) {
		// P2 == 0: return P1
		x3.Set(x1)
		y3.Set(y1)
		z3.Set(z1)
		return
	}

	// z1² and z2²
	z1_pow2 := new(big.Int).Mul(z1, z1)
	mod(z1_pow2)
	z2_pow2 := new(big.Int).Mul(z2, z2)
	mod(z2_pow2)

	// u1 = x1 * z2²
	u1 := new(big.Int).Mul(x1, z2_pow2)
	mod(u1)

	// u2 = x2 * z1²
	u2 := new(big.Int).Mul(x2, z1_pow2)
	mod(u2)

	// z1³
	z1_pow3 := z1_pow2.Mul(z1_pow2, z1)
	z1_pow2 = nil
	mod(z1_pow3)

	// z2³
	z2_pow3 := z2_pow2.Mul(z2_pow2, z2)
	z2_pow2 = nil
	mod(z2_pow3)

	// s1 = y1 * z2³
	s1 := z2_pow3.Mul(y1, z2_pow3)
	z2_pow3 = nil
	mod(s1)

	// s2 = y2 * z1³
	s2 := z1_pow3.Mul(y2, z1_pow3)
	z1_pow3 = nil
	mod(s2)

	// h = u2 - u1
	h := u2.Sub(u2, u1)
	u2 = nil
	mod(h)

	// r = s2 - s1
	r := s2.Sub(s2, s1)
	s2 = nil
	mod(r)

	//  h = (x2 * z1²) - (x1 * z2²)
	//  r = (y2 * z1³) - (y1 * z2³)
	//
	// Affine coordinates are calculated as:
	//  Ax1 = x1 / z1²
	//  Ax2 = x2 / z2²
	//  Ay1 = y1 / z1³
	//  Ay2 = y2 / z2³
	//
	// Thus, if h = 0, the X-coordinate is the same:
	//  x2 * z1² = x1 * z2²
	//  x2 / z2² = x1 / z1²
	//  Ax1 = Ax2
	//
	// and, if r = 0, the Y-coordinate is the same:
	//  y2 * z1³ = y1 * z2³
	//  y2 / z2³ = y1 / z1³
	//  Ay1 = Ay2
	if equal(h, zero) {
		if equal(r, zero) {
			// P1 == P2: return the doubled point
			DoubleJacobi(
				x1, y1, z1,
				x3, y3, z3,
			)
			return
		}

		// P1 == -P2: sum will be zero
		// INVARIANT: for performance, y2 is assumed to be negative y1
		x3.Set(zero)
		y3.Set(zero)
		z3.Set(zero)
		return
	}

	// h²
	hh := new(big.Int).Mul(h, h)
	mod(hh)

	// v = u1 * h²
	v := u1.Mul(u1, hh)
	u1 = nil
	mod(v)

	// h³
	hhh := hh.Mul(hh, h)
	hh = nil
	mod(hhh)

	// x3 = r² - h³ - 2*v
	x3.Mul(r, r)
	x3.Sub(x3, hhh)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	mod(x3)

	// y3 = r * (v - x3) - s1 * h³
	y3.Sub(v, x3)
	y3.Mul(r, y3)
	y3.Sub(y3, s1.Mul(s1, hhh))
	s1 = nil
	mod(y3)

	// z3 = z1 * z2 * h
	z3.Mul(z1, z2)
	z3.Mul(z3, h)
	mod(z3)
}

// AddAffine adds two affine points together:
//
//  P1 + P2 = P3
//  (x1, y1) + (x2, y2) = (x3, y3)
//
// It stores the resulting affine point in x3, and y3.
//
// We incorporate the standard affine addition and doubling formulas:
//
//  if P1 == P2:
//   m = (3 * x1² + a) / (2 * y1)
//  else:
//   m = (y2 - y1) / (x2 - x1)
//  x3 = m² - x1 - x2
//  y3 = m * (x1 - x3) - y1
//
// This function does not check point validity - it assumes you
// are passing valid points on the secp256k1 curve.
func AddAffine(
	x1, y1 *big.Int,
	x2, y2 *big.Int,
	x3, y3 *big.Int,
) {
	// P2 + 0 = P2
	if EqualAffine(x1, y1, zero, zero) {
		x3.Set(x2)
		y3.Set(y2)
		return
	}
	// P1 + 0 = P1
	if EqualAffine(x2, y2, zero, zero) {
		x3.Set(x1)
		y3.Set(y1)
		return
	}

	xEqual := equal(x1, x2)
	yEqual := equal(y1, y2)

	// INVARIANT: if x1 == x2 && y1 != y2, assume y1 = -y2 (the only other possible point on the curve).
	// Thus P1 + P2 = 0
	if xEqual && !yEqual {
		x3.Set(zero)
		y3.Set(zero)
		return
	}

	m := new(big.Int)

	if xEqual && yEqual {
		// m = (3 * x1² + a) / (2 * y1)
		m.Mul(x1, x1)
		m.Mul(m, three)
		twoY1Inverse := x3.Mul(y1, two)
		invert(twoY1Inverse)
		m.Mul(m, twoY1Inverse)
	} else {
		//  m = (y2 - y1) / (x2 - x1)
		m.Sub(y2, y1)
		xDiffInverse := x3.Sub(x2, x1)
		invert(xDiffInverse)
		m.Mul(m, xDiffInverse)
	}
	mod(m)

	// x3 = m² - x1 - x2
	x3.Mul(m, m)
	x3.Sub(x3, x1)
	x3.Sub(x3, x2)
	mod(x3)

	// y3 = m * (x1 - x3) - y1
	y3.Sub(x1, x3)
	y3.Mul(y3, m)
	y3.Sub(y3, y1)
	mod(y3)
}
