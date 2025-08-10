package main

import (
	"fmt"
	"strings"

	"github.com/mmcloughlin/avo/attr"
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	"github.com/quant1x/x/labs/asm/avx2"
)

// 单精度浮点数组a和b计算, 结果保存在result, 返回剩余多少个元素(arithmetic,binary operation)
func float32x8_binary_template(prefix, operator string) {
	prefix = strings.TrimSpace(prefix)
	operator = strings.TrimSpace(operator)
	function := fmt.Sprintf("%s_%s", prefix, operator)

	TEXT(function, attr.NOSPLIT, "func(a, b, result []float32) int")
	Pragma("noescape")
	Doc(fmt.Sprintf("%s %ss a and b, and store the results in result.", function, operator))
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
	case avx2.Op_add:
		VADDPS(Y1, Y0, Y0)
	case avx2.Op_sub:
		VSUBPS(Y1, Y0, Y0)
	case avx2.Op_mul:
		VMULPS(Y1, Y0, Y0)
	case avx2.Op_div:
		VDIVPS(Y1, Y0, Y0)
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

//// 生成 float32切片 加 的plan9汇编代码
//func genFloat32x8Add() {
//	TEXT("f32x8_add", attr.NOSPLIT, "func(a, b, result []float32) int")
//	Pragma("noescape")
//	a := Mem{Base: Load(Param("a").Base(), GP64())}
//	b := Mem{Base: Load(Param("b").Base(), GP64())}
//	c := Mem{Base: Load(Param("result").Base(), GP64())}
//	n := Load(Param("a").Len(), GP64())
//
//	Y0 := YMM()
//	Y1 := YMM()
//
//	Label("loop")
//	CMPQ(n, U32(8))
//	JL(LabelRef("done"))
//
//	VMOVUPS(a.Offset(0), Y0)
//	VMOVUPS(b.Offset(0), Y1)
//	VADDPS(Y1, Y0, Y0)
//	VMOVUPS(Y0, c.Offset(0))
//
//	ADDQ(U32(32), a.Base)
//	ADDQ(U32(32), b.Base)
//	ADDQ(U32(32), c.Base)
//	SUBQ(U32(8), n)
//	JMP(LabelRef("loop"))
//
//	Label("done")
//	Store(n, ReturnIndex(0))
//	VZEROUPPER()
//	RET()
//
//	//Generate()
//}
//
//// 生成 float32切片 减 的plan9汇编代码
//func genFloat32x8Sub() {
//	TEXT("f32x8_sub", attr.NOSPLIT, "func(a, b, result []float32) int")
//	Pragma("noescape")
//	a := Mem{Base: Load(Param("a").Base(), GP64())}
//	b := Mem{Base: Load(Param("b").Base(), GP64())}
//	c := Mem{Base: Load(Param("result").Base(), GP64())}
//	n := Load(Param("a").Len(), GP64())
//
//	Y0 := YMM()
//	Y1 := YMM()
//
//	Label("loop")
//	CMPQ(n, U32(8))
//	JL(LabelRef("done"))
//
//	VMOVUPS(a.Offset(0), Y0)
//	VMOVUPS(b.Offset(0), Y1)
//	VSUBPS(Y1, Y0, Y0)
//	VMOVUPS(Y0, c.Offset(0))
//
//	ADDQ(U32(32), a.Base)
//	ADDQ(U32(32), b.Base)
//	ADDQ(U32(32), c.Base)
//	SUBQ(U32(8), n)
//	JMP(LabelRef("loop"))
//
//	Label("done")
//	Store(n, ReturnIndex(0))
//	VZEROUPPER()
//	RET()
//
//	//Generate()
//}
//
//// 生成 float32切片 乘 的plan9汇编代码
//func genFloat32x8Mul() {
//	TEXT("f32x8_mul", attr.NOSPLIT, "func(a, b, result []float32) int")
//	Pragma("noescape")
//	a := Mem{Base: Load(Param("a").Base(), GP64())}
//	b := Mem{Base: Load(Param("b").Base(), GP64())}
//	c := Mem{Base: Load(Param("result").Base(), GP64())}
//	n := Load(Param("a").Len(), GP64())
//
//	Y0 := YMM()
//	Y1 := YMM()
//
//	Label("loop")
//	CMPQ(n, U32(8))
//	JL(LabelRef("done"))
//
//	VMOVUPS(a.Offset(0), Y0)
//	VMOVUPS(b.Offset(0), Y1)
//	VSUBPS(Y1, Y0, Y0)
//	VMOVUPS(Y0, c.Offset(0))
//
//	ADDQ(U32(32), a.Base)
//	ADDQ(U32(32), b.Base)
//	ADDQ(U32(32), c.Base)
//	SUBQ(U32(8), n)
//	JMP(LabelRef("loop"))
//
//	Label("done")
//	Store(n, ReturnIndex(0))
//	VZEROUPPER()
//	RET()
//
//	//Generate()
//}
