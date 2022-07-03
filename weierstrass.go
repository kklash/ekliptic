package ekliptic

import (
	"math/big"
)

// Weierstrass solves the Weierstrass form elliptic curve equation for y, given an affine value of x:
//  y² = x³ + ax + b
// ...where a = 0 and b = 7 - these are secp256k1 curve constants.
//
// It returns the two possible values of y: an even value and an odd value. You could also call this function "GetYValues(x)".
// This is useful for uncompressing public keys, which in secp256k1 often supply only the-x coordinate, and maybe a specifier
// indicating whether the y-coordinate is even or odd.
//
// Returns two nil values if x is not a valid coordinate on the secp256k1 curve.
func Weierstrass(x *big.Int) (evenY, oddY *big.Int) {
	if x.Cmp(Secp256k1_P) >= 0 {
		// x is larger than curve modulus P; point would not be on the curve
		return nil, nil
	}

	// c = x³ + ax + b
	c := new(big.Int)
	c.Mul(x, x)
	c.Mul(c, x)
	c.Add(c, Secp256k1_B)
	mod(c)

	// this is actually faster than using big.Int's ModSqrt method.
	y := new(big.Int).Exp(c, squareRootExp, Secp256k1_P) // y = c^((p+1)/4)

	ySquared := new(big.Int).Mul(y, y)
	mod(ySquared)
	if !equal(c, ySquared) {
		// c != y² mod p - this means the given x-coordinate is not on the curve
		return nil, nil
	}

	// if y is even, -y is odd, and vice-versa.
	evenY = y
	oddY = ySquared.Set(y)
	Negate(oddY)

	y = nil
	ySquared = nil

	if !isEven(evenY) {
		evenY, oddY = oddY, evenY
	}

	return evenY, oddY
}

func isEven(y *big.Int) bool {
	return y.Bit(0) == 0
}
