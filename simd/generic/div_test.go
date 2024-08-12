package generic

import (
	"math"
	"testing"
)

func Test_div(t *testing.T) {
	type args struct {
		a float32
		b float32
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			name: "div-0",
			args: args{a: 0, b: 0},
			want: float32(math.NaN()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := div(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("div() = %v, want %v", got, tt.want)
			}
		})
	}
}
