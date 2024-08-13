package internal

import (
	"fmt"
	"testing"
)

func Test__mm256_mul_ps1(t *testing.T) {
	a := Float32x8{3, 3, 3, 3, 3, 3, 3, 3}
	b := Float32x8{2, 2, 2, 2, 2, 2, 2, 2}
	n := len(a)
	fmt.Println(n)
	fmt.Printf("a: %p\n", &a)
	fmt.Printf("b: %p\n", &b)
	fmt.Printf("n: %p\n", &n)
	//ptrA := unsafe.Pointer(&a[0])
	//ptrB := unsafe.Pointer(&b[0])
	//ptrC := unsafe.Pointer(&c[0])
	c := _mm256_mul_ps1(a, b)
	fmt.Println(c)
}
