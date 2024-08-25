package internal

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test_avx2_mm256_float32_add(t *testing.T) {
	a := Float32x8{3, 3, 3, 3, 3, 3, 3, 3}
	b := Float32x8{2, 2, 2, 2, 2, 2, 2, 2}
	n := len(a)
	c := make([]Float32x8, n)
	fmt.Println(n)
	fmt.Printf("a: %p\n", &a)
	fmt.Printf("b: %p\n", &b)
	fmt.Printf("n: %p, %x\n", &n, unsafe.Pointer(uintptr(n)))
	ptrA := unsafe.Pointer(&a[0])
	ptrB := unsafe.Pointer(&b[0])
	ptrC := unsafe.Pointer(&c[0])
	_avx2_mm256_float32_add(ptrA, ptrB, ptrC, unsafe.Pointer(uintptr(n)))
	fmt.Println(c)
}

//func Test_inline_asm(t *testing.T) {
//	var value int64 = 42
//	var result int64
//
//	asm(
//		"XORQ %0, %0", // 汇编指令：将寄存器自身与自身进行 XOR 操作，结果为 0
//		"MOVQ %0, %1", // 将结果从 %0 寄存器移动到变量 result 中
//		"NOP",         // 无操作，仅用于演示
//		&value,        // 第一个输入/输出寄存器，使用地址
//		&result,       // 第二个输入/输出寄存器，使用地址
//		"",            // 无 clobber 列表
//		"volatile",    // 指示编译器此汇编代码是易变的，不应被优化掉
//	)
//
//	fmt.Println("Result:", result) // 应该输出 0
//}
