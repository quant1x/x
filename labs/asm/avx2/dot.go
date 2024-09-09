package avx2

func DotNaive(a, b []float32) float32 {
	sum := float32(0)
	for i := 0; i < len(a) && i < len(b); i++ {
		sum += a[i] * b[i]
	}
	return sum
}
