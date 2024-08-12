package simd

import (
	"reflect"
	"testing"
)

func TestFloat32x4_Div(t *testing.T) {
	type args struct {
		other Float32x4
	}
	tests := []struct {
		name    string
		m       Float32x4
		args    args
		want    Float32x4
		wantErr bool
	}{
		{
			name: "test1",
			m:    Float32x4{1, 2, 3, 4},
			args: args{other: Float32x4{1, 2, 3, 4}},
			want: Float32x4{1, 1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Div(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Div() got = %v, want %v", got, tt.want)
			}
		})
	}
}
