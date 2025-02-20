package climate

import (
	"gabe565.com/ambient-weather-fusion/pkg/constraints"
)

// CtoF converts Celsius to Fahrenheit.
func CtoF[V constraints.Number](c V) float64 {
	return float64(c)*9/5 + 32
}

// FtoC converts Fahrenheit to Celsius.
func FtoC[V constraints.Number](f V) float64 {
	return (float64(f) - 32) * 5 / 9
}

const kmhToMPHConversionFactor = 0.621371192

// KMHtoMPH converts kilometers-per hour to miles-per-hour.
func KMHtoMPH[V constraints.Number](kmh V) float64 {
	return float64(kmh) * kmhToMPHConversionFactor
}

// MPHtoKMH converts miles-per-hour to kilometers-per-hour.
func MPHtoKMH[V constraints.Number](mph V) float64 {
	return float64(mph) / kmhToMPHConversionFactor
}
