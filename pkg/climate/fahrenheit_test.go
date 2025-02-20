package climate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDewPointF(t *testing.T) {
	type args struct {
		tempF    float64
		humidity float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"50F at 70%", args{50, 70}, 40.60648803127103},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, DewPointF(tt.args.tempF, tt.args.humidity), 0.000001)
		})
	}
}

func TestWindChillF(t *testing.T) {
	type args struct {
		tempF        float64
		windSpeedMPH float64
	}
	type testCase struct {
		name string
		args args
		want float64
	}
	tests := []testCase{
		{"32F at 0mph", args{32, 0}, 32},
		{"32F at 3mph", args{32, 3}, 29.316734684752944},
		{"50F at 10mph", args{50, 10}, 46.03680329552729},
		{"86F at 10mph", args{86, 10}, 86},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, WindChillF(tt.args.tempF, tt.args.windSpeedMPH), 0.000001)
		})
	}
}

func TestHeatIndexF(t *testing.T) {
	type args struct {
		tempF    float64
		humidity float64
	}
	type testCase struct {
		name string
		args args
		want float64
	}
	tests := []testCase{
		{"69F at 70%", args{69, 70}, 68.89},
		{"86F at 80%", args{80, 80}, 84.23041600000002},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, HeatIndexF(tt.args.tempF, tt.args.humidity), 0.000001)
		})
	}
}

func TestFeelsLikeF(t *testing.T) {
	type args struct {
		tempF        float64
		humidity     float64
		windSpeedMPH float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"50F, 70%, 10mph", args{50, 70, 10}, 46.03680329552729},
		{"70F, 70%, 10mph", args{70, 70, 10}, 69.99000000000001},
		{"40F, 70%, 10mph", args{40, 70, 10}, 33.64254827558847},
		{"90F, 70%, 10mph", args{90, 70, 10}, 105.92202060000027},
		{"90F, 10%, 10mph", args{90, 10, 10}, 85.27896836218746},
		{"80F, 90%, 10mph", args{80, 90, 10}, 86.34189169999989},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feelsLike := FeelsLikeF(tt.args.tempF, tt.args.humidity, tt.args.windSpeedMPH)
			assert.InDelta(t, tt.want, feelsLike, 0.000001)
		})
	}
}
