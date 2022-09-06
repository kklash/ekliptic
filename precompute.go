package ekliptic

import "math/big"

//go:generate go run ./genprecompute -o precomputed_table.go

// PrecomputedTable is a table in the form of a 2d slice, holding precomputed
// multiplications of a point P. Each entry in the table is the product of P and
// some multiple of a power of 2.
//
// The table is laid out like this:
//
//        0 1 2 3 4 ... 15
//  -----------------------
//  2^0  |* * * * * ...  *
//  2^4  |* * * * * ...  *
//  2^8  |* * * * * ...  *
//  2^16 |* * * * * ...  *
//  ...  |* * * * * ...  *
//  2^252|* * * * * ...  *
//
// Each row i of the table holds the multiplications of the point, doubled 4i times.
// i should be less than 64.
//
// Each column j of the table holds the doublings of the point, multiplied j times.
// j should be less than 16.
//
// Indexing the table like so:
//
//  table[i][j]
//
// ...Looks up the following precomputed point multiplication:
//
//  2^(4i) * j * P
type PrecomputedTable [][][2]*big.Int

// NewPrecomputedTable computes a PrecomputedTable used for speeding up multiplications
// of a fixed affine point (x, y).
func NewPrecomputedTable(x, y *big.Int) PrecomputedTable {
	table := make(PrecomputedTable, 64)
	for i := range table {
		table[i] = make([][2]*big.Int, 16)
		table[i][0] = [2]*big.Int{new(big.Int), new(big.Int)}
		for j := 1; j < len(table[i]); j++ {
			table[i][j][0], table[i][j][1] = AddAffine(
				table[i][j-1][0], table[i][j-1][1],
				x, y,
			)
		}
		if i+1 < len(table) {
			for k := 0; k < 4; k++ {
				x, y = DoubleAffine(x, y)
			}
		}
	}
	return table
}
