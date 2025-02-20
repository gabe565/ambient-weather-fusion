package geolocation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const statueLatitude, statueLongitude = 40.6892, -74.0445

func statueOfLiberty() Point { return Point{statueLatitude, statueLongitude} }

func assertPointInDelta(t *testing.T, want, got Point) {
	assert.InDelta(t, want.Latitude, got.Latitude, 0.000001)
	assert.InDelta(t, want.Longitude, got.Longitude, 0.000001)
}

func TestPoint_div(t *testing.T) {
	type args struct {
		k float64
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Point
	}{
		{"statue of liberty", statueOfLiberty(), args{2}, Pt(statueLatitude/2, statueLongitude/2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPointInDelta(t, tt.want, tt.p.div(tt.args.k))
		})
	}
}

func TestPoint_mul(t *testing.T) {
	type args struct {
		k float64
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Point
	}{
		{"statue of liberty", statueOfLiberty(), args{2}, Pt(statueLatitude*2, statueLongitude*2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPointInDelta(t, tt.want, tt.p.mul(tt.args.k))
		})
	}
}

func TestPoint_ShiftPoint(t *testing.T) {
	type args struct {
		delta Point
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Point
	}{
		{
			"statue of liberty shift 4",
			statueOfLiberty(),
			args{Pt(4, 4)},
			Pt(40.74708914323121, -73.96815500763032),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPointInDelta(t, tt.want, tt.p.ShiftPoint(tt.args.delta))
		})
	}
}

func TestPoint_String(t *testing.T) {
	tests := []struct {
		name string
		p    Point
		want string
	}{
		{"statue of liberty", statueOfLiberty(), "(40.6892, -74.0445)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.String())
		})
	}
}

func TestPt(t *testing.T) {
	type args struct {
		latitude  float64
		longitude float64
	}
	tests := []struct {
		name string
		args args
		want Point
	}{
		{"statue of liberty", args{statueLatitude, statueLongitude}, statueOfLiberty()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPointInDelta(t, tt.want, Pt(tt.args.latitude, tt.args.longitude))
		})
	}
}
