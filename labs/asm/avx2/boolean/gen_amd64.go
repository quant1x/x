package main

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/quant1x/x/labs/asm/avx2"
)

func main() {
	bool_binary_template(avx2.Func_prefix_boolx8, avx2.Op_and)
	bool_binary_template(avx2.Func_prefix_boolx8, avx2.Op_or)
	bool_binary_template(avx2.Func_prefix_boolx8, avx2.Op_xor)
	bool_i8x32_binary_template(avx2.Func_prefix_boolx32, avx2.Op_and)
	build.Generate()
}
