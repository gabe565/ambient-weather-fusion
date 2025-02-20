package climate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtoF(t *testing.T) {
	type args struct {
		c float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"freezing", args{0}, 32},
		{"boiling", args{100}, 212},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, CtoF(tt.args.c), 0.000001)
		})
	}
}

func TestFtoC(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"freezing", args{32}, 0},
		{"boiling", args{212}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, FtoC(tt.args.f), 0.000001)
		})
	}
}

func TestKMHtoMPH(t *testing.T) {
	type args struct {
		mph float64
	}
	type testCase struct {
		name string
		args args
		want float64
	}
	tests := []testCase{
		{"0kmh", args{0}, 0},
		{"10kmh", args{16.09344000614692}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, KMHtoMPH(tt.args.mph), 0.000001)
		})
	}
}

func TestMPHtoKMH(t *testing.T) {
	type args struct {
		mph float64
	}
	type testCase struct {
		name string
		args args
		want float64
	}
	tests := []testCase{
		{"0mph", args{0}, 0},
		{"10mph", args{10}, 16.09344000614692},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, MPHtoKMH(tt.args.mph), 0.000001)
		})
	}
}
