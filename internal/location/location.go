package location

import "math"

const earthRadius = 3963.1

func Shift(lat, lon, latDelta, lonDelta float64) (float64, float64) {
	// Convert to radians
	latRad := lat * math.Pi / 180
	lonRad := lon * math.Pi / 180

	// Calculate angular distance in radians
	angularLatDelta := latDelta / earthRadius
	angularLonDelta := lonDelta / earthRadius

	// Calculate new lat/lon
	newLatRad := latRad + angularLatDelta
	newLonRad := lonRad + angularLonDelta/math.Cos(latRad)

	// Convert new lat/lon to degrees
	shiftedLat := newLatRad * 180 / math.Pi
	shiftedLon := newLonRad * 180 / math.Pi

	return shiftedLat, shiftedLon
}
