package ekliptic

import "testing"

func BenchmarkNewPrecomputedTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPrecomputedTable(Secp256k1_GeneratorX, Secp256k1_GeneratorY)
	}
}
