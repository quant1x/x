package sse2

// sse2 - 128b

// type Boolx16 [16]bool
// type Boolx8 [8]bool
// type Boolx4 [4]bool
// type Boolx2 [2]bool

type Int8x16 [16]int8
type Int16x8 [8]int16
type Int32x4 [4]int32
type Int64x2 [2]int64

type Uint8x16 [16]uint8
type Uint16x8 [8]uint16
type Uint32x4 [4]uint32
type Uint64x2 [2]uint64

type Float32x4 [4]float32
type Float64x2 [2]float64

func (f *Float32x4) from(x []float32) {
	sn := len(x) % 4
	for i := 0; i < sn; i++ {
		f[i] = x[i]
	}
}

func (f *Float32x4) as_array() []float32 {
	return f[:]
}
