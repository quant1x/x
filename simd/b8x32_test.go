package simd

import (
	"slices"
	"testing"
)

// AND 正式版本
func Benchmark_bools_and_release(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	//r := make([]bool, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bs_and(x, y)
	}
}

func Benchmark_bools_and_noasm(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	//r := make([]bool, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = noasm_b8x32_and(x, y)
	}
}

func Benchmark_bools_and_v1(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1_bs_and(x, y)
	}
}

func Benchmark_bools_and_v2(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v2_bs_and(x, y)
	}
}

// for vek
func Benchmark_bools_and_v3_for_vek(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v3_bs_and(x, y)
	}
}

// OR 正式版本
func Benchmark_bools_or_release(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	//r := make([]bool, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bs_or(x, y)
	}
}

func Benchmark_bools_or_noasm(b *testing.B) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	//r := make([]bool, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = noasm_b8x32_or(x, y)
	}
}
