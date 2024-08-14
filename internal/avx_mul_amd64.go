package internal

import "unsafe"

//go:noescape
func avx2_mm256_float32_add(a, b, c, n unsafe.Pointer)
