package climate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDewPointC(t *testing.T) {
	type args struct {
		tempC    float64
		humidity float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"10C at 70%", args{10, 70}, 4.781382239595014},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, DewPointC(tt.args.tempC, tt.args.humidity), 0.000001)
		})
	}
}

func TestWindChillC(t *testing.T) {
	type args struct {
		tempC        float64
		windSpeedKMH float64
	}
	type testCase struct {
		name string
		args args
		want float64
	}
	tests := []testCase{
		{"0C at 0kmh", args{0, 0}, 0},
		{"0C at 5kmh", args{0, 5}, -1.572787454594585},
		{"10C at 16kmh", args{10, 16}, 7.8089738644784825},
		{"30C at 16kmh", args{30, 16}, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, WindChillC(tt.args.tempC, tt.args.windSpeedKMH), 0.000001)
		})
	}
}

func TestHeatIndexC(t *testing.T) {
	type args struct {
		tempC    float64
		humidity float64
	}
	type testCase struct {
		name string
		args args
		want float64
	}
	tests := []testCase{
		{"15C at 70%", args{15, 70}, 14.383333333333333},
		{"21C at 70%", args{21, 70}, 20.983333333333338},
		{"30C at 80%", args{30, 80}, 37.66703727777778},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, HeatIndexC(tt.args.tempC, tt.args.humidity), 0.000001)
		})
	}
}

func TestFeelsLikeC(t *testing.T) {
	type args struct {
		tempC        float64
		humidity     float64
		windSpeedKPH float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"50F, 70%, 10mph", args{26.6667, 90, 16.0934}, 30.190028154626233},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feelsLike := FeelsLikeC(tt.args.tempC, tt.args.humidity, tt.args.windSpeedKPH)
			assert.InDelta(t, tt.want, feelsLike, 0.000001)
		})
	}
}
