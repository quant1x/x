package simd

import (
	"slices"
	"testing"
)

func Benchmark_bools_bs_and(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	//r := make([]bool, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bs_and(x, y)
	}
}

func Benchmark_bools_bs_and_v1(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1_bs_and(x, y)
	}
}

func Benchmark_bools_bs_and_v2(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v2_bs_and(x, y)
	}
}

func Benchmark_bs_and_noasm(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	//r := make([]bool, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = noasm_b8x32_and(x, y)
	}
}
