package geolocation

import (
	"math"
	"strconv"
)

// EarthRadius is the radius of the Earth.
const EarthRadius = 3959.0

// Pt is shorthand for Point{Latitude, Longitude}.
func Pt(latitude, longitude float64) Point {
	return Point{
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// Point represents a geographic coordinate with a latitude and longitude.
type Point struct {
	Latitude, Longitude float64
}

// String returns a string representation of p like "(40.6892, -74.0445)".
func (p Point) String() string {
	return "(" + strconv.FormatFloat(p.Latitude, 'f', -1, 64) +
		", " + strconv.FormatFloat(p.Longitude, 'f', -1, 64) + ")"
}

func (p Point) mul(k float64) Point {
	return Point{Latitude: p.Latitude * k, Longitude: p.Longitude * k}
}

func (p Point) div(k float64) Point {
	return Point{Latitude: p.Latitude / k, Longitude: p.Longitude / k}
}

func (p Point) Radians() Point {
	return p.mul(math.Pi / 180)
}

func (p Point) Degrees() Point {
	return p.mul(180 / math.Pi)
}

// ShiftPoint returns a new Point that is shifted by the given delta in miles.
// The delta's Latitude value specifies the north/south displacement,
// and the delta's Longitude value specifies the east/west displacement.
// The shift roughly takes into account the curvature of the Earth by converting miles to an angular distance.
func (p Point) ShiftPoint(q Point) Point {
	radians := p.Radians()
	angularDelta := q.div(EarthRadius)
	return Pt(
		radians.Latitude+angularDelta.Latitude,
		radians.Longitude+angularDelta.Longitude/math.Cos(radians.Latitude),
	).Degrees()
}

// Shift is shorthand for Point.ShiftPoint(Pt(Latitude, Longitude)).
func (p Point) Shift(latDelta, longDelta float64) Point {
	return p.ShiftPoint(Pt(latDelta, longDelta))
}
