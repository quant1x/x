package simd

import (
	"fmt"
	"testing"

	"github.com/quant1x/x/assert"
)

func Benchmark_mm256_add_ps(tb *testing.B) {
	size := 1000
	a := make([]float32, size)
	b := make([]float32, size)
	c := make([]float32, size)

	for i := 0; i < size; i++ {
		a[i] = float32(i)
		b[i] = float32(i)
	}

	tb.ResetTimer()
	for i := 0; i < tb.N; i++ {
		_mm256_add_ps(a, b, c)
	}
}

func Test_avx2_float32_add(t *testing.T) {
	a := []float32{1, 2, 3, 4, 5, 6, 7, 8}
	b := []float32{2, 4, 6, 8, 10, 12, 14, 16}
	ok := []float32{3, 6, 9, 12, 15, 18, 21, 24}
	//r := make([]float32, len(a))
	r := avx2_float32_add(a, b)
	if !assert.Equal(r, ok) {
		t.Error("错误", r)
	}
}

func Test_avx2_float32_add_v2(t *testing.T) {
	a := []float32{1, 2, 3, 4, 5, 6, 7}
	b := []float32{2, 4, 6, 8, 10, 12, 14}
	ok := []float32{0, 0, 0, 0, 0, 0, 0}
	//r := make([]float32, len(a))
	r := avx2_float32_add(a, b)
	if !assert.Equal(r, ok) {
		t.Error("错误", r)
	}
}

func Test_avx2_f32_add(t *testing.T) {
	a := f32x8{1, 2, 3, 4, 5, 6, 7, 8}
	b := f32x8{2, 4, 6, 8, 10, 12, 14, 16}
	ok := f32x8{3, 6, 9, 12, 15, 18, 21, 24}
	r := avx2_f32x8_add(a, b)
	if r != ok {
		t.Error("错误", r)
	}
	fmt.Println(r)
}

func Benchmark_avx2_f32_add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := f32x8{1, 2, 3, 4, 5, 6, 7, 8}
		y := f32x8{2, 4, 6, 8, 10, 12, 14, 16}
		ok := f32x8{3, 6, 9, 12, 15, 18, 21, 24}
		r := avx2_f32x8_add(x, y)
		if r != ok {
			b.Error("错误", r)
		}
	}
}
