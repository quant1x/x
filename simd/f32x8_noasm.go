package simd

func noasm_f32x8_add(x, y []float32) []float32 {
	n := len(x)
	r := make([]float32, n)
	for i := 0; i < n; i++ {
		r[i] = x[i] + y[i]
	}
	return r
}
