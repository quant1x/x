package simd

import "unsafe"

type AutoPtr unsafe.Pointer

type register int

const (
	sse2   = 128
	avx2   = 256
	avx512 = 512
)

type SimdPtr[E number] struct {
	s     *[]E
	lanes int
	bits  int
}

func (p *SimdPtr[E]) pointer() unsafe.Pointer {
	return unsafe.Pointer(&(*p.s)[0])
}

func (p *SimdPtr[E]) firstAddress() uintptr {
	return uintptr(p.pointer())
}

func (p *SimdPtr[E]) from(s []E) {
	p.s = &s
	p.lanes = 8
}

func (p *SimdPtr[E]) seek(n int) unsafe.Pointer {
	ptr := p.firstAddress()
	ptr += uintptr(n)
	return unsafe.Pointer(ptr)
}

func (p *SimdPtr[E]) asInt8s() []int8 {
	return (*[32]int8)(p.pointer())[:]
}

func (p *SimdPtr[E]) asInt16s() []int16 {
	return (*[16]int16)(p.pointer())[:]
}

func (p *SimdPtr[E]) asInt32s() []int32 {
	return (*[8]int32)(p.pointer())[:]
}

func (p *SimdPtr[E]) asInt64s() []int64 {
	return (*[4]int64)(p.pointer())[:]
}
