package simd

import (
	"fmt"
	"testing"
	"unsafe"
)

func test_gen_int8s() []int8 {
	length := 32
	ret := make([]int8, length)
	for i := 0; i < length; i++ {
		ret[i] = int8(i)
	}
	return ret
}

func Test_load(t *testing.T) {
	var s = test_gen_int8s()
	fmt.Printf("s address: %p\n", &s)
	var p SimdPtr[int8]
	p.from(s)
	fmt.Printf("p address: %p\n", p.pointer())
	for i := 0; i < len(s); i++ {
		fmt.Printf("s address[%d]: %p\n", i, &s[i])
	}
	var d = load(unsafe.Pointer(&s[0]))
	fmt.Printf("d address: %p\n", &d)
	fmt.Println(s)
	fmt.Println(d)
}
