package main

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/quant1x/x/labs/asm/avx2"
)

func main() {
	float32x8_binary_template(avx2.Func_prefix_float32x8, avx2.Op_add)
	float32x8_binary_template(avx2.Func_prefix_float32x8, avx2.Op_sub)
	float32x8_binary_template(avx2.Func_prefix_float32x8, avx2.Op_mul)
	float32x8_binary_template(avx2.Func_prefix_float32x8, avx2.Op_div)
	float32x8_binary_template(avx2.Func_prefix_float32x8, avx2.Op_and)
	float32x8_binary_template(avx2.Func_prefix_float32x8, avx2.Op_or)
	float32x8_binary_template(avx2.Func_prefix_float32x8, avx2.Op_xor)
	build.Generate()
}
