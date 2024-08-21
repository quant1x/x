package simd

import (
	"golang.org/x/sys/cpu"
)

var (
	__has_sse2 = false
	__has_avx  = false
	__has_avx2 = false
)

func init() {
	__has_sse2 = cpu.X86.HasSSE2
	__has_avx = cpu.X86.HasAVX
	__has_avx2 = cpu.X86.HasAVX2
}

//func Add(a, b, result []float32) {
//	an := len(a)
//	bn := len(b)
//	if an != bn {
//		panic("Add: bad len")
//	}
//	length := an
//	// for avx2
//	epoch := length / 8
//	remain := length % 8
//	start := 0
//	if __has_avx2 && epoch > 0 {
//		for i := 0; i < epoch; i++ {
//			pos := start + i*8
//			var x, y f32x8
//			x.from(a, pos)
//			y.from(b, pos)
//			r := avx2_f32x8_add(x, y)
//			result = append(result, r.as_array()...)
//		}
//	}
//	start += epoch * 8
//	// for sse2
//	epoch = remain / 4
//	remain = remain % 4
//	if __has_sse2 && epoch > 0 {
//		for i := 0; i < epoch; i++ {
//			pos := start + i*4
//			var x, y sse2.Float32x4
//			x.from(a[pos : pos+4])
//			y.from(b[pos : pos+4])
//			r := sse2.AddFloat32x4(x, y)
//			result = append(result, r.as_array()...)
//		}
//	}
//	start += epoch * 4
//	// for common
//	if remain > 0 {
//		r := noasm_f32x8_add(a[start:], b[start:])
//		result = append(result, r...)
//	}
//}
