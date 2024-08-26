package simd

import (
	"fmt"
	"testing"
)

func TestSimdPtr_asInt64s(t *testing.T) {
	var s = test_gen_int8s()
	fmt.Printf("s address: %p\n", &s)
	var p1 SimdPtr[int8]
	p1.from(s)
	fmt.Printf("p1 address: %p\n", p1.pointer())
	int64s := p1.asInt64s()
	fmt.Printf("int64s address: %p, %v\n", &int64s, int64s)
	fmt.Printf("0: %x, 1: %x, 2: %x, 3: %x\n", int64s[0], int64s[1], int64s[2], int64s[3])

	var p2 SimdPtr[int64]
	p2.from(int64s)
	fmt.Printf("p2 address: %p\n", p2.pointer())
	s2 := p2.asInt8s()
	fmt.Printf("s2 address: %p\n", &s2)
	fmt.Printf("s2: %v\n", s2)
	p2.seek(1)
}
