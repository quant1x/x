package simd

import (
	"fmt"
	"testing"
)

func Test_f32x8_v3from(t *testing.T) {
	a := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	var b f32x8
	b.v3from(a, 0)
	fmt.Println(b)
	b.v3from(a, 8)
	fmt.Println(b)
}
