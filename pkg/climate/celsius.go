package climate

import (
	"math"

	"gabe565.com/ambient-weather-fusion/pkg/constraints"
)

const (
	magnusA = 17.27
	magnusB = 237.7
)

// DewPointC computes the dew point in Celsius.
func DewPointC[Temp, Humidity constraints.Number](tempC Temp, humidity Humidity) float64 {
	if humidity > 100 {
		humidity = 100
	}
	g := magnusA*float64(tempC)/(magnusB+float64(tempC)) + math.Log(float64(humidity)/100.0)
	return magnusB * g / (magnusA - g)
}

// WindChillC computes the wind chill in Celsius.
func WindChillC[Temp, Humidity constraints.Number](tempC Temp, windSpeedKMH Humidity) float64 {
	tempF := CtoF(tempC)
	windSpeedMPH := KMHtoMPH(windSpeedKMH)
	windChillF := WindChillF(tempF, windSpeedMPH)
	return FtoC(windChillF)
}

// HeatIndexC computes the heat index in Celsius.
func HeatIndexC[Temp, Humidity constraints.Number](tempC Temp, humidity Humidity) float64 {
	tempF := CtoF(tempC)
	heatIndexF := HeatIndexF(tempF, humidity)
	return FtoC(heatIndexF)
}

// FeelsLikeC computes the feels-like temperature in Celsius.
func FeelsLikeC[Temp, Humidity, WindSpeed constraints.Number](
	tempC Temp,
	humidity Humidity,
	windSpeedKMH WindSpeed,
) float64 {
	tempF := CtoF(tempC)
	windSpeedMPH := KMHtoMPH(windSpeedKMH)
	feelsLikeF := FeelsLikeF(tempF, humidity, windSpeedMPH)
	return FtoC(feelsLikeF)
}
