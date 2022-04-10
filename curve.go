package ekliptic

import (
	"crypto/elliptic"
	"math/big"
)

// Curve satisfies crypto/elliptic.Curve using the secp256k1 curve paramters.
type Curve struct {
	params *elliptic.CurveParams
}

// Params returns the parameters for the curve. Satisfies elliptic.Curve.
func (c *Curve) Params() *elliptic.CurveParams {
	if c.params == nil {
		c.params = &elliptic.CurveParams{
			P:       new(big.Int).Set(Secp256k1_P),
			N:       new(big.Int).Set(Secp256k1_CurveOrder),
			B:       new(big.Int).Set(Secp256k1_B),
			Gx:      new(big.Int).Set(Secp256k1_GeneratorX),
			Gy:      new(big.Int).Set(Secp256k1_GeneratorY),
			BitSize: 256,
			Name:    "secp256k1",
		}
	}
	return c.params
}

// IsOnCurve reports whether the given (x,y) lies on the curve. Satisfies elliptic.Curve.
// Note: The elliptic.Curve interface requires that infinity is NOT on the curve.
func (_ *Curve) IsOnCurve(x, y *big.Int) bool {
	if equal(x, zero) && equal(y, zero) {
		return false
	}
	return IsOnCurveAffine(x, y)
}

// Add returns the sum of (x1,y1) and (x2,y2) satisfies elliptic.Curve.
func (_ *Curve) Add(x1, y1, x2, y2 *big.Int) (x3, y3 *big.Int) {
	x3 = new(big.Int)
	y3 = new(big.Int)
	AddAffine(x1, y1, x2, y2, x3, y3)
	return
}

// Double returns 2*(x,y). Satisfies elliptic.Curve.
func (_ *Curve) Double(x1, y1 *big.Int) (x3, y3 *big.Int) {
	x3 = new(big.Int)
	y3 = new(big.Int)
	DoubleAffine(x1, y1, x3, y3)
	return
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
// Satisfies elliptic.Curve.
func (_ *Curve) ScalarMult(x1, y1 *big.Int, k []byte) (x2, y2 *big.Int) {
	x2 = new(big.Int)
	y2 = new(big.Int)
	kBig := new(big.Int).SetBytes(k)

	if equal(x1, Secp256k1_GeneratorX) && equal(y1, Secp256k1_GeneratorY) {
		MultiplyBasePoint(kBig, x2, y2)
	} else {
		MultiplyAffine(x1, y1, kBig, x2, y2, nil)
	}
	return
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form. Satisfies elliptic.Curve.
func (_ *Curve) ScalarBaseMult(k []byte) (x2, y2 *big.Int) {
	x2 = new(big.Int)
	y2 = new(big.Int)
	kBig := new(big.Int).SetBytes(k)
	MultiplyBasePoint(kBig, x2, y2)
	return
}
