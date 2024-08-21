package simd

//go:noescape
func avx2_f32x8_add(x, y f32x8) f32x8

// func _mm256_add_ps(x, y, z unsafe.Pointer)
//
//go:noescape
func _mm256_add_ps(a []float32, b []float32, c []float32) int

//go:noescape
func Add_AVX2_F32(x []float32, y []float32)
