package ekliptic

import (
	"math/big"
)

// Weierstrass solves the Weierstrass form elliptic curve equation for y,
// given an affine value of x:
//
//	y² = x³ + ax + b mod P
//
// ...where a = 0, b = 7 and P is Secp256k1_P - the secp256k1 curve constants.
//
// It returns the two possible values of y: an even value and an odd value. You could also
// call this function "GetYValues(x)". This is useful for uncompressing public keys, which
// in secp256k1 often supply only the-x coordinate, and maybe a specifier indicating whether
// the y-coordinate is even or odd.
//
// Returns two nil values if x is not a valid coordinate on the secp256k1 curve.
// Returns two zero values if x is zero.
func Weierstrass(x *big.Int) (evenY, oddY *big.Int) {
	// if x is not positive, or is larger than curve modulus P, point would not be on the curve.
	if equal(x, zero) {
		return new(big.Int).Set(zero), new(big.Int).Set(zero)
	} else if x.Cmp(zero) == -1 || x.Cmp(Secp256k1_P) >= 0 {
		return nil, nil
	}

	// c = x³ + ax + b
	c := new(big.Int)
	c.Mul(x, x)
	c.Mul(c, x)
	c.Add(c, Secp256k1_B)
	modCoordinate(c)

	// this is actually faster than using big.Int's ModSqrt method.
	y := new(big.Int).Exp(c, squareRootExp, Secp256k1_P) // y = c^((p+1)/4)

	ySquared := new(big.Int).Mul(y, y)
	modCoordinate(ySquared)
	if !equal(c, ySquared) {
		// c != y² mod p - this means the given x-coordinate is not on the curve
		return nil, nil
	}

	// if y is even, -y is odd, and vice-versa.
	evenY = y
	oddY = Negate(y)

	if !isEven(evenY) {
		evenY, oddY = oddY, evenY
	}

	return evenY, oddY
}

func isEven(y *big.Int) bool {
	return y.Bit(0) == 0
}
