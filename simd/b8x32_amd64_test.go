package simd

import (
	"slices"
	"testing"

	"github.com/quant1x/x/assert"
)

func Test_b1x8_and(t *testing.T) {
	type args struct {
		a      []bool
		b      []bool
		result []bool
		want   []bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_count_7",
			args: args{
				a:      []bool{true, true, true, true, true, true, true},
				b:      []bool{true, true, true, true, true, true, true},
				result: make([]bool, 7),
				want:   []bool{false, false, false, false, false, false, false},
			},
			want: 7,
		},
		{
			name: "test_count_8",
			args: args{
				a:      []bool{false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
				b:      []bool{true, true, true, true, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
				result: make([]bool, 32),
				want:   []bool{false, true, true, true, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b32x8_and(tt.args.a, tt.args.b, tt.args.result)
			if got != tt.want || !assert.Equal(tt.args.want, tt.args.result) {
				t.Errorf("b1x8_and() = %v, want %v, result = %v, want = %v", got, tt.want, tt.args.result, tt.args.want)
			}
		})
	}
}

func Test_b8x32_and(t *testing.T) {
	type args struct {
		a      []bool
		b      []bool
		result []bool
		want   []bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_count_7",
			args: args{
				a:      []bool{true, true, true, true, true, true, true},
				b:      []bool{true, true, true, true, true, true, true},
				result: make([]bool, 7),
				want:   []bool{false, false, false, false, false, false, false},
			},
			want: 7,
		},
		{
			name: "test_count_8",
			args: args{
				a:      []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, true},
				b:      []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
				result: make([]bool, 32),
				want:   []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, true},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b8x32_and(tt.args.a, tt.args.b, tt.args.result)
			if got != tt.want || !assert.Equal(tt.args.want, tt.args.result) {
				t.Errorf("b8x32_and() = %v, want %v, result = %v, want = %v", got, tt.want, tt.args.result, tt.args.want)
			}
		})
	}
}

func Test_b8x32_and_check(t *testing.T) {
	x := slices.Clone(testDataBoolx)
	y := slices.Clone(testDataBooly)
	got := bs_and(x, y)
	if !assert.Equal(got, testDataBoolr) {
		t.Errorf("bs_and() = %v, want %v", got, testDataBoolr)
	}
}
