package simd

import (
	"github.com/quant1x/x/assert"
	"testing"
)

func Test_f32x8_add(t *testing.T) {
	type args struct {
		a      []float32
		b      []float32
		result []float32
		want   []float32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_count_7",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7},
				b:      []float32{1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 7),
				want:   []float32{0, 0, 0, 0, 0, 0, 0},
			},
			want: 7,
		},
		{
			name: "test_count_8",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7, 8},
				b:      []float32{1, 1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 8),
				want:   []float32{2, 3, 4, 5, 6, 7, 8, 9},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f32x8_add(tt.args.a, tt.args.b, tt.args.result)
			if got != tt.want || !assert.Equal(tt.args.want, tt.args.result) {
				t.Errorf("f32x8_add() = %v, want %v, result = %v, want = %v", got, tt.want, tt.args.result, tt.args.want)
			}
		})
	}
}

func Test_f32x8_sub(t *testing.T) {
	type args struct {
		a      []float32
		b      []float32
		result []float32
		want   []float32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_count_7",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7},
				b:      []float32{1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 7),
				want:   []float32{0, 0, 0, 0, 0, 0, 0},
			},
			want: 7,
		},
		{
			name: "test_count_8",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7, 8},
				b:      []float32{1, 1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 8),
				want:   []float32{0, 1, 2, 3, 4, 5, 6, 7},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f32x8_sub(tt.args.a, tt.args.b, tt.args.result)
			if got != tt.want || !assert.Equal(tt.args.want, tt.args.result) {
				t.Errorf("f32x8_sub() = %v, want %v, result = %v, want = %v", got, tt.want, tt.args.result, tt.args.want)
			}
		})
	}
}

func Test_f32x8_mul(t *testing.T) {
	type args struct {
		a      []float32
		b      []float32
		result []float32
		want   []float32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_count_7",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7},
				b:      []float32{1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 7),
				want:   []float32{0, 0, 0, 0, 0, 0, 0},
			},
			want: 7,
		},
		{
			name: "test_count_8",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7, 8},
				b:      []float32{1, 1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 8),
				want:   []float32{1, 2, 3, 4, 5, 6, 7, 8},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f32x8_mul(tt.args.a, tt.args.b, tt.args.result)
			if got != tt.want || !assert.Equal(tt.args.want, tt.args.result) {
				t.Errorf("f32x8_mul() = %v, want %v, result = %v, want = %v", got, tt.want, tt.args.result, tt.args.want)
			}
		})
	}
}

func Test_f32x8_div(t *testing.T) {
	type args struct {
		a      []float32
		b      []float32
		result []float32
		want   []float32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_count_7",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7},
				b:      []float32{1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 7),
				want:   []float32{0, 0, 0, 0, 0, 0, 0},
			},
			want: 7,
		},
		{
			name: "test_count_8",
			args: args{
				a:      []float32{1, 2, 3, 4, 5, 6, 7, 8},
				b:      []float32{1, 1, 1, 1, 1, 1, 1, 1},
				result: make([]float32, 8),
				want:   []float32{1, 2, 3, 4, 5, 6, 7, 8},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f32x8_div(tt.args.a, tt.args.b, tt.args.result)
			if got != tt.want || !assert.Equal(tt.args.want, tt.args.result) {
				t.Errorf("f32x8_div() = %v, want %v, result = %v, want = %v", got, tt.want, tt.args.result, tt.args.want)
			}
		})
	}
}

//func Test_f32x8_and(t *testing.T) {
//	type args struct {
//		a      []float32
//		b      []float32
//		result []float32
//		want   []float32
//	}
//	tests := []struct {
//		name string
//		args args
//		want int
//	}{
//		{
//			name: "test_count_7",
//			args: args{
//				a:      []float32{1, 2, 3, 4, 5, 6, 7},
//				b:      []float32{1, 1, 1, 1, 1, 1, 1},
//				result: make([]float32, 7),
//				want:   []float32{0, 0, 0, 0, 0, 0, 0},
//			},
//			want: 7,
//		},
//		{
//			name: "test_count_8",
//			args: args{
//				a:      []float32{1, 2, 3, 4, 5, 6, 7, 8},
//				b:      []float32{1, 1, 1, 1, 1, 1, 1, 1},
//				result: make([]float32, 8),
//				want:   []float32{1, 0, 3, 4, 5, 6, 7, 8},
//			},
//			want: 0,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := f32x8_and(tt.args.a, tt.args.b, tt.args.result)
//			if got != tt.want || !assert.Equal(tt.args.want, tt.args.result) {
//				t.Errorf("f32x8_and() = %v, want %v, result = %v, want = %v", got, tt.want, tt.args.result, tt.args.want)
//			}
//		})
//	}
//}
