package simd

import (
	"unsafe"

	"github.com/quant1x/x/std"
)

func convect[E1 number, E2 number](from []E1) []E2 {
	length1 := len(from)
	first1 := from[0]
	size1 := int(unsafe.Sizeof(first1))
	//ptr1 := unsafe.Pointer(&first1)
	ptr1 := unsafe.SliceData(from)
	var tmpE2 E2
	length2 := length1 * size1 / int(unsafe.Sizeof(tmpE2))
	ptr2 := (*E2)(unsafe.Pointer(ptr1))
	return unsafe.Slice(ptr2, length2)
}

type AutoPtr unsafe.Pointer

type register int

const (
	sse2   = 128
	avx2   = 256
	avx512 = 512
)

type Pointer[E number] struct {
	s      *[]E
	lanes  int
	bits   int
	length int
}

func (p *Pointer[E]) pointer() unsafe.Pointer {
	return unsafe.Pointer(&(*p.s)[0])
}

func (p *Pointer[E]) firstAddress() uintptr {
	return uintptr(p.pointer())
}

func (p *Pointer[E]) from(s []E) {
	p.length = len(s)
	p.s = &s
	p.bits = 8 * int(unsafe.Sizeof((*p.s)[0]))
	p.lanes = avx2 / p.bits
}

func (p *Pointer[E]) seek(n int) unsafe.Pointer {
	//ptr := p.firstAddress()
	//ptr += uintptr(n)
	ptr := p.pointer()
	return unsafe.Pointer(uintptr(ptr) + uintptr(n))
}

func (p *Pointer[E]) asInt8s() []int8 {
	return (*[32]int8)(p.pointer())[:]
}

func (p *Pointer[E]) asInt16s() []int16 {
	return (*[16]int16)(p.pointer())[:]
}

func (p *Pointer[E]) asInt32s() []int32 {
	return (*[8]int32)(p.pointer())[:]
}

func (p *Pointer[E]) asInt64s() []int64 {
	return (*[4]int64)(p.pointer())[:]
}

func (p *Pointer[E]) asBools() []bool {
	return (*[32]bool)(p.pointer())[:]
}

func (p *Pointer[E]) HexString() string {
	length := len(*p.s)
	data := unsafe.Slice((*byte)(p.pointer()), length*p.bits/8)
	return std.ToHexString(data)
}
