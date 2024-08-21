package simd

import (
	"github.com/quant1x/x/assert"
	"slices"
	"testing"
)

func Benchmark_add_none(b *testing.B) {
	x := slices.Clone(testDataFloat32)
	y := slices.Clone(testDataFloat32y)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = noasm_f32x8_add(x, y)
	}
}

func Benchmark_add_v1(b *testing.B) {
	x := slices.Clone(testDataFloat32)
	y := slices.Clone(testDataFloat32y)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = avx2_float32_add(x, y)
	}
}

func Benchmark_add_v2(b *testing.B) {
	x := slices.Clone(testDataFloat32)
	y := slices.Clone(testDataFloat32y)
	r := make([]float32, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f32x8_add(x, y, r)
	}
}

func Benchmark_add_v3(b *testing.B) {
	x := slices.Clone(testDataFloat32)
	y := slices.Clone(testDataFloat32y)
	//// 开始创建公共arena内存池
	//mem := arena.NewArena()
	//// 最后统一释放内存池
	//defer mem.Free()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := f32s_add(x, y)
		clear(r)
	}
}

func Benchmark_add_vek32(b *testing.B) {
	x := slices.Clone(testDataFloat32)
	y := slices.Clone(testDataFloat32y)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vek_f32x8_Add(x, y)
	}
}

func Test_f32s_add(t *testing.T) {
	type args struct {
		a []float32
		b []float32
		//result []float32
		//want   []float32
	}
	tests := []struct {
		name string
		args args
		want []float32
	}{
		{
			name: "test_count_7",
			args: args{
				a: []float32{1, 2, 3, 4, 5, 6, 7},
				b: []float32{1, 1, 1, 1, 1, 1, 1},
			},
			want: []float32{2, 3, 4, 5, 6, 7, 8},
		},
		{
			name: "test_count_8",
			args: args{
				a: []float32{1, 2, 3, 4, 5, 6, 7, 8},
				b: []float32{1, 1, 1, 1, 1, 1, 1, 1},
			},
			want: []float32{2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f32s_add(tt.args.a, tt.args.b)
			if !assert.Equal(got, tt.want) {
				t.Errorf("f32s_add() = %v, want %v", got, tt.want)
			}
		})
	}
}
