package simd

import "unsafe"

type i8x32 [32]int8

//go:noescape
func load(p unsafe.Pointer) i8x32
