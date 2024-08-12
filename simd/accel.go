package simd

func (m Float32x4) Div(other Float32x4) (Float32x4, error) {
	return DivFloat32x4(m, other)
}
