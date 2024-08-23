package simd

func noasm_b8x32_and(x, y []bool) []bool {
	n := len(x)
	r := make([]bool, n)
	for i := 0; i < n; i++ {
		r[i] = x[i] && y[i]
	}
	return r
}

func noasm_b8x32_or(x, y []bool) []bool {
	n := len(x)
	r := make([]bool, n)
	for i := 0; i < n; i++ {
		r[i] = x[i] || y[i]
	}
	return r
}

func noasm_b8x32_xor(x, y []bool) []bool {
	n := len(x)
	r := make([]bool, n)
	for i := 0; i < n; i++ {
		r[i] = bool_xor(x[i], y[i])
	}
	return r
}
