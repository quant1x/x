package simd

func f32s_add(a, b []float32) []float32 {
	an := len(a)
	bn := len(b)
	if an != bn {
		panic("Add: bad len")
	}
	length := an
	result := make([]float32, length, length)
	n := f32x8_add(a, b, result)

	for i := length - n; i < length; i++ {
		result[i] = a[i] + b[i]
	}
	return result
}
