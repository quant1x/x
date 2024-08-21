package main

import (
	"fmt"
	"github.com/mmcloughlin/avo/attr"
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	"github.com/quant1x/x/labs/asm/avx2"
	"strings"
)

// 布尔数组a和b逻辑计算, 结果保存在result, 返回剩余多少个元素(arithmetic,binary operation)
func bool_binary_template(prefix, operator string) {
	prefix = strings.TrimSpace(prefix)
	operator = strings.TrimSpace(operator)
	function := fmt.Sprintf("%s_%s", prefix, operator)

	TEXT(function, attr.NOSPLIT, "func(a, b, result []bool) int")
	Pragma("noescape")
	a := Mem{Base: Load(Param("a").Base(), GP64())}
	b := Mem{Base: Load(Param("b").Base(), GP64())}
	c := Mem{Base: Load(Param("result").Base(), GP64())}
	n := Load(Param("a").Len(), GP64())

	Y0 := YMM()
	Y1 := YMM()

	Label("loop")
	CMPQ(n, U32(8))
	JL(LabelRef("done"))

	VMOVUPS(a.Offset(0), Y0)
	VMOVUPS(b.Offset(0), Y1)
	switch operator {
	//case avx2.Op_add:
	//	VADDPS(Y1, Y0, Y0)
	//case avx2.Op_sub:
	//	VSUBPS(Y1, Y0, Y0)
	//case avx2.Op_mul:
	//	VMULPS(Y1, Y0, Y0)
	//case avx2.Op_div:
	//	VDIVPS(Y1, Y0, Y0)
	case avx2.Op_and:
		VANDPS(Y1, Y0, Y0)
	case avx2.Op_or:
		VORPS(Y1, Y0, Y0)
	case avx2.Op_xor:
		VXORPS(Y1, Y0, Y0)
	default:
		panic("not implemented: " + operator)
	}
	VMOVUPS(Y0, c.Offset(0))

	ADDQ(U32(32), a.Base)
	ADDQ(U32(32), b.Base)
	ADDQ(U32(32), c.Base)
	SUBQ(U32(8), n)
	JMP(LabelRef("loop"))

	Label("done")
	Store(n, ReturnIndex(0))
	VZEROUPPER()
	RET()

	//Generate()
}

// 布尔数组a和b逻辑计算, 结果保存在result, 返回剩余多少个元素(arithmetic,binary operation)
func bool_i8x32_binary_template(prefix, operator string) {
	prefix = strings.TrimSpace(prefix)
	operator = strings.TrimSpace(operator)
	function := fmt.Sprintf("%s_%s", prefix, operator)

	TEXT(function, attr.NOSPLIT, "func(a, b, result []bool) int")
	Pragma("noescape")
	a := Mem{Base: Load(Param("a").Base(), GP64())}
	b := Mem{Base: Load(Param("b").Base(), GP64())}
	c := Mem{Base: Load(Param("result").Base(), GP64())}
	n := Load(Param("a").Len(), GP64())

	Y0 := YMM()
	Y1 := YMM()

	Label("loop")
	CMPQ(n, U32(32))
	JL(LabelRef("done"))

	VMOVDQU(a.Offset(0), Y0)
	VMOVDQU(b.Offset(0), Y1)
	switch operator {
	//case avx2.Op_add:
	//	VADDPS(Y1, Y0, Y0)
	//case avx2.Op_sub:
	//	VSUBPS(Y1, Y0, Y0)
	//case avx2.Op_mul:
	//	VMULPS(Y1, Y0, Y0)
	//case avx2.Op_div:
	//	VDIVPS(Y1, Y0, Y0)
	case avx2.Op_and:
		//VANDPS(Y1, Y0, Y0)
		VPAND(Y1, Y0, Y0)
	case avx2.Op_or:
		VORPS(Y1, Y0, Y0)
	case avx2.Op_xor:
		VXORPS(Y1, Y0, Y0)
	default:
		panic("not implemented: " + operator)
	}
	VMOVDQU(Y0, c.Offset(0))

	ADDQ(U32(32), a.Base)
	ADDQ(U32(32), b.Base)
	ADDQ(U32(32), c.Base)
	SUBQ(U32(32), n)
	JMP(LabelRef("loop"))

	Label("done")
	Store(n, ReturnIndex(0))
	VZEROUPPER()
	RET()

	//Generate()
}
