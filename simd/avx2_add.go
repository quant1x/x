package simd

import "slices"

// Add returns the result of adding two slices element-wise.
func vek_f32x8_Add(x, y []float32) []float32 {
	x = slices.Clone(x)
	Add_Inplace(x, y)
	return x
}

// Add_Inplace adds a slice element-wise to the first slice, inplace.
func Add_Inplace(x, y []float32) {
	Add_AVX2_F32(x, y)
}

//func _float32_add(x, y, r []float32) {
//	ptrA := unsafe.Pointer(&x[0])
//	ptrB := unsafe.Pointer(&y[0])
//	ptrC := unsafe.Pointer(&r[0])
//	_mm256_add_ps(ptrA, ptrB, ptrC)
//}

func avx2_float32_add(a, b []float32) []float32 {
	//lot_size := 8
	//an := len(a)
	//bn := len(b)
	//if an != bn {
	//	panic("Add: bad len")
	//}
	//length := an
	//result := make([]float32, length)
	//// for avx2
	//epoch := length / 8
	//remain := length % 8
	//start := 0
	//for i := 0; i < epoch; i++ {
	//	pos := start + i*8
	//	_float32_add(a[pos:], b[pos:], result[pos:])
	//}
	//start += epoch * 8
	//for i := 0; i < remain; i++ {
	//	pos := start + i
	//	result[pos] = a[pos] + b[pos]
	//}
	//
	//return result
	an := len(a)
	bn := len(b)
	if an != bn {
		panic("Add: bad len")
	}
	length := an
	result := make([]float32, length)
	n := _mm256_add_ps(a, b, result)
	//fmt.Println(n)
	_ = n
	return result
}

//func avx2_float32_add(x, y []float32) []float32 {
//	n := len(x)
//	r := make([]float32, n)
//	ptrA := unsafe.Pointer(&x[0])
//	ptrB := unsafe.Pointer(&y[0])
//	ptrC := unsafe.Pointer(&r[0])
//	ptrN := unsafe.Pointer(uintptr(n))
//	avx2_mm256_float32_add(ptrA, ptrB, ptrC, ptrN)
//	return r
//}
