package ekliptic

import (
	"fmt"
	"math/big"
	"testing"
	"time"
)

const multOps = 100_000

func TestMultiplyNonModulo(t *testing.T) {
	t.Skip()

	x := hexint("262136d6f958d318d187fc7bde50c31d8cec89b388203ae3cee4b1d6fb367fd2")
	y := hexint("0e899411dad3b1d2f91e5b406d828f87809dcaa955d3c0461aa3f932ecf8dc4d")
	z := new(big.Int)
	n := new(big.Int)

	for i := 0; i < 40; i++ {
		start := time.Now()
		for j := 0; j < multOps; j++ {
			z.Mul(x, y)
		}
		nBits := len(z.Bytes()) * 8
		fmt.Printf("%d-bits x 256-bits: %d ns/op\n", nBits, (int(time.Since(start)) / multOps))
		x.Set(z)

		start = time.Now()
		for j := 0; j < multOps; j++ {
			n.Set(z)
			modCoordinate(n)
		}

		fmt.Printf("mod(%d-bits): %d ns/op\n\n", nBits, (int(time.Since(start)) / multOps))
	}
}
