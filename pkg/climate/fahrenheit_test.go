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
	t.Run("boundary", func(t *testing.T) {
		tests := []testCase{
			{"temp>50 returns raw", args{86, 10}, 86},
			{"wind<3 returns raw", args{32, 2.9}, 32},
			{"temp=51 returns raw", args{51, 10}, 51},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.InDelta(t, tt.want, WindChillF(tt.args.tempF, tt.args.windSpeedMPH), 0.000001)
			})
		}
	})

	// NWS Wind Chill Chart reference values (rounded to nearest integer).
	// Source: https://www.weather.gov/safety/cold-wind-chill-chart
	t.Run("NWS chart", func(t *testing.T) {
		tests := []testCase{
			{"40F at 10mph", args{40, 10}, 34},
			{"30F at 15mph", args{30, 15}, 19},
			{"25F at 25mph", args{25, 25}, 9},
			{"0F at 15mph", args{0, 15}, -19},
			{"10F at 40mph", args{10, 40}, -15},
			{"-20F at 30mph", args{-20, 30}, -53},
			{"-40F at 60mph", args{-40, 60}, -91},
			{"5F at 5mph", args{5, 5}, -5},
			{"20F at 20mph", args{20, 20}, 4},
			{"-10F at 45mph", args{-10, 45}, -44},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.InDelta(t, tt.want, WindChillF(tt.args.tempF, tt.args.windSpeedMPH), 0.5)
			})
		}
	})
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
	t.Run("below 80F", func(t *testing.T) {
		// Simple formula path (temp < 80), no NWS chart coverage.
		tests := []testCase{
			{"69F at 70%", args{69, 70}, 68.89},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.InDelta(t, tt.want, HeatIndexF(tt.args.tempF, tt.args.humidity), 0.000001)
			})
		}
	})

	// NWS Heat Index Chart reference values (rounded to nearest integer).
	// Source: https://www.weather.gov/safety/heat-index
	t.Run("NWS chart", func(t *testing.T) {
		tests := []testCase{
			// Base Rothfusz regression
			{"80F at 40%", args{80, 40}, 80},
			{"90F at 40%", args{90, 40}, 91},
			{"86F at 50%", args{86, 50}, 88},
			{"96F at 50%", args{96, 50}, 108},
			{"90F at 70%", args{90, 70}, 105},
			{"100F at 60%", args{100, 60}, 129},
			// High humidity adjustment (humidity > 85, temp <= 87)
			{"80F at 85%", args{80, 85}, 85},
			{"82F at 90%", args{82, 90}, 91},
			{"86F at 90%", args{86, 90}, 105},
			// Various coverage
			{"94F at 75%", args{94, 75}, 124},
			{"104F at 40%", args{104, 40}, 119},
			{"88F at 80%", args{88, 80}, 106},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.InDelta(t, tt.want, HeatIndexF(tt.args.tempF, tt.args.humidity), 1)
			})
		}
	})
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
		{"boundary 50F with wind", args{50, 50, 10}, 46.03680329552729},
		{"boundary 68F no effect", args{68, 50, 5}, 68},
		{"boundary 69F heat index", args{69, 50, 5}, 67.94999999999999},
		{"dead zone 55F", args{55, 50, 10}, 55},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feelsLike := FeelsLikeF(tt.args.tempF, tt.args.humidity, tt.args.windSpeedMPH)
			assert.InDelta(t, tt.want, feelsLike, 0.000001)
		})
	}
}
