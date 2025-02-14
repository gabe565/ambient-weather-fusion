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
		{"temperature 0", args{0, 70}, 0},
		{"humidity 0", args{10, 0}, 0},
		{"10C at 70%", args{10, 70}, 4.781382239595014},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, DewPointC(tt.args.tempC, tt.args.humidity), 0.000001)
		})
	}
}

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
		{"temperature 0", args{0, 70}, 0},
		{"humidity 0", args{10, 0}, 0},
		{"50F at 70%", args{50, 70}, 40.60648803127103},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, DewPointF(tt.args.tempF, tt.args.humidity), 0.000001)
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
		{"temperature 0", args{0, 70, 10}, 0},
		{"humidity 0", args{50, 0, 10}, 0},
		{"wind speed 0", args{50, 70, 0}, 0},
		{"50F, 70%, 10mph", args{50, 70, 10}, 50.0},
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

func TestFeelsLikeC(t *testing.T) {
	type args struct {
		tempF        float64
		humidity     float64
		windSpeedKPH float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"temperature 0", args{0, 70, 16}, 0},
		{"humidity 0", args{26, 0, 16}, 0},
		{"wind speed 0", args{50, 70, 0}, 0},
		{"50F, 70%, 10mph", args{26.6667, 90, 16.0934}, 30.190028154626233},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feelsLike := FeelsLikeC(tt.args.tempF, tt.args.humidity, tt.args.windSpeedKPH)
			if tt.want == 0 {
				assert.InDelta(t, tt.want, feelsLike, 0)
			} else {
				assert.InEpsilon(t, tt.want, feelsLike, 0.000001)
			}
		})
	}
}
