package avx

import (
	"fmt"
	"github.com/quant1x/x/simd"
	"math"
	"testing"
	"unsafe"
)

func Test_mm256_mul_ps(t *testing.T) {
	//a := []float32{1, 1, 1, 1, 1, 1, 1, 1}
	a := []float32{3, 3, 3, 3, 3, 3, 3, float32(math.NaN())}
	b := []float32{2, 2, 2, 2, 2, 2, 2, 0.0}
	c := make([]float32, len(a))
	n := len(a)
	fmt.Println(n)
	fmt.Printf("a: %p\n", &a)
	fmt.Printf("b: %p\n", &b)
	fmt.Printf("c: %p\n", &c)
	fmt.Printf("n: %p\n", &n)
	ptrA := unsafe.Pointer(&a[0])
	ptrB := unsafe.Pointer(&b[0])
	ptrC := unsafe.Pointer(&c[0])
	mm256_mul_ps(ptrA, ptrB, ptrC)
	fmt.Println(c)
}

func Test_mm256_mul_ps_v2(t *testing.T) {
	//a := []float32{1, 1, 1, 1, 1, 1, 1, 1}
	a := simd.Float32x8{3, 3, 3, 3, 3, 3, 3, 3}
	b := simd.Float32x8{2, 2, 2, 2, 2, 2, 2, 2}
	c := make([]float32, len(a))
	n := len(a)
	fmt.Println(n)
	fmt.Printf("a: %p\n", &a)
	fmt.Printf("b: %p\n", &b)
	fmt.Printf("c: %p\n", &c)
	fmt.Printf("n: %p\n", &n)
	ptrA := unsafe.Pointer(&a[0])
	ptrB := unsafe.Pointer(&b[0])
	ptrC := unsafe.Pointer(&c[0])
	mm256_mul_ps(ptrA, ptrB, ptrC)
	fmt.Println(c)
}
