package ekliptic

import "math/big"

// DoubleJacobi doubles a Jacobian coordinate point on the secp256k1 curve, using the "dbl-2009-l" doubling formulas.
//
// http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
//  A = X1²
//  B = Y1²
//  C = B²
//  D = 2*((X1+B)²-A-C)
//  E = 3*A
//  F = E²
//  X3 = F-2*D
//  Y3 = E*(D-X3)-8*C
//  Z3 = 2*Y1*Z1
func DoubleJacobi(
	x1, y1, z1 *big.Int,
	x3, y3, z3 *big.Int,
) {
	// a = x1²
	a := new(big.Int).Mul(x1, x1)

	// b = y1²
	b := new(big.Int).Mul(y1, y1)

	// c = b²
	c := new(big.Int).Mul(b, b)

	// The official dbl-2009-l formula specifies:
	//  d = 2 * ((x1+b)² - a - c)
	//
	// But it can be simplified, because:
	//  a = x1², c = b²
	//  d = 2 * ((x1+b)² - a - c)
	//    = 2 * ((x1+b)² - x1² - b²)
	//    = 2 * ((x1+b)(x1+b) - x1² - b²)
	//    = 2 * (x1² + 2(x1*b) + b² - x1² - b²)
	//    = 2 * 2(x1*b)
	//
	// So actually:
	//  d = 4 * x1 * b
	d := b.Mul(b, x1)
	d.Mul(d, four)
	b = nil

	// e = 3 * a
	e := a.Mul(a, three)
	a = nil

	// f = e²
	f := new(big.Int).Mul(e, e)

	// x3 = f - 2 * d
	x3.Mul(d, two)
	x3.Sub(f, x3)
	mod(x3)

	// z3 = 2 * y1 * z1
	z3.Mul(y1, z1)
	z3.Mul(z3, two)
	mod(z3)

	// *** Ensure y3 is set AFTER z3. If y3 points to the same bigint  ***
	// *** as y1, this will mutate y1, which is needed to calculate z3 ***

	// y3 = e * (d - x3) - 8 * c
	y3.Sub(d, x3)
	y3.Mul(y3, e)
	y3.Sub(y3, c.Mul(c, eight))
	c = nil
	mod(y3)
}

//  m = (3*x1²+a) / (2*y1)
//  x3 = m² - x1 - x2
//  y3 = m(x1-x3) - y1
func DoubleAffine(
	x1, y1 *big.Int,
	x3, y3 *big.Int,
) {
	AddAffine(x1, y1, x1, y1, x3, y3)
}
