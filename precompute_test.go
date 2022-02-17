package ekliptic

import "testing"

func BenchmarkComputePointDoubles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ComputePointDoubles(Secp256k1_GeneratorX, Secp256k1_GeneratorY)
	}
}
