package simd

import "unsafe"

type AutoPtr unsafe.Pointer

type SimdPtr[E number] struct {
	s *[]E
}

func (p *SimdPtr[E]) pointer() unsafe.Pointer {
	return unsafe.Pointer(&(*p.s)[0])
}

func (p *SimdPtr[E]) firstAddress() uintptr {
	return uintptr(p.pointer())
}

func (p *SimdPtr[E]) load(s []E) {
	p.s = &s
}

func (p *SimdPtr[E]) offset(n int) unsafe.Pointer {
	ptr := p.firstAddress()
	ptr += uintptr(n)
	return unsafe.Pointer(ptr)
}

func (p *SimdPtr[E]) asInt8s() []int8 {
	return (*[32]int8)(p.pointer())[:]
}

func (p *SimdPtr[E]) asInt64s() []int64 {
	return (*[4]int64)(p.pointer())[:]
}
