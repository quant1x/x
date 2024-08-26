package simd

func noasm_f32x8_add(x, y []float32) []float32 {
	n := len(x)
	r := make([]float32, n)
	for i := 0; i < n; i++ {
		r[i] = x[i] + y[i]
	}
	return r
}

func noasm_f32x8_add_v1(x, y, r []float32) int {
	length := len(x)
	n := length
	step := 1
	for i := 0; i < length; i += step {
		r[i] = x[i] + y[i]
		n -= step
	}
	return n
}
