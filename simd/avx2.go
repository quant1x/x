package simd

import "unsafe"

var (
	sizeOfFloat32 = unsafe.Sizeof(float32(0))
)

// avx2 - 256b

type Int8x32 [32]int8
type Int16x16 [16]int16
type Int32x8 [8]int32
type Int64x4 [4]int64

type Uint8x32 [32]uint8
type Uint16x16 [16]uint16
type Uint32x8 [8]uint32
type Uint64x4 [4]uint64

type f32x8 [8]float32
type f64x4 [4]float64

func (f *f32x8) from(x []float32, offset int) {
	f.v3from(x, offset)
}
func (f *f32x8) v1from(x []float32, offset int) {
	sn := len(x) % 8
	for i := 0; i < sn; i++ {
		f[i] = x[i]
	}
}

func (f *f32x8) v2from(x []float32, offset int) {
	*f = f32x8(x)
}
func (f *f32x8) v3from(x []float32, offset int) {
	length := 8
	ptr := unsafe.Pointer(&x[0])
	addr := uintptr(ptr)
	addr += uintptr(offset) * sizeOfFloat32
	e := unsafe.Slice((*float32)(unsafe.Pointer(addr)), length)
	*f = [8]float32(e)
}

func (f *f32x8) as_array() []float32 {
	//return f[:]
	length := 8
	ptr := unsafe.Pointer(&f[0])
	addr := uintptr(ptr)
	e := unsafe.Slice((*float32)(unsafe.Pointer(addr)), length)
	return e
}
