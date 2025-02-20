package climate

import (
	"math"

	"gabe565.com/ambient-weather-fusion/pkg/constraints"
)

// DewPointF computes the dew point in Fahrenheit.
func DewPointF[Temp, Humidity constraints.Number](tempF Temp, humidity Humidity) float64 {
	tempC := FtoC(tempF)
	dewPoint := DewPointC(tempC, humidity)
	return CtoF(dewPoint)
}

// WindChillF computes the wind chill in Fahrenheit.
func WindChillF[Temp, Humidity constraints.Number](tempF Temp, windSpeedMPH Humidity) float64 {
	if windSpeedMPH < 3 || tempF > 50 {
		return float64(tempF)
	}
	exp := math.Pow(float64(windSpeedMPH), 0.16)
	return 35.74 + 0.6215*float64(tempF) - 35.75*exp + 0.4275*float64(tempF)*exp
}

// HeatIndexF computes the heat index in Fahrenheit.
func HeatIndexF[Temp, Humidity constraints.Number](tempF Temp, humidity Humidity) float64 {
	if tempF < 80 {
		return 0.5 * (float64(tempF) + 61 + (float64(tempF)-68)*1.2 + float64(humidity)*0.094)
	}

	base := -42.379 +
		2.04901523*float64(tempF) +
		10.14333127*float64(humidity) +
		-0.22475541*float64(tempF)*float64(humidity) +
		-0.00683783*float64(tempF)*float64(tempF) +
		-0.05481717*float64(humidity)*float64(humidity) +
		0.00122874*float64(tempF)*float64(tempF)*float64(humidity) +
		0.00085282*float64(tempF)*float64(humidity)*float64(humidity) +
		-0.00000199*float64(tempF)*float64(tempF)*float64(humidity)*float64(humidity)
	switch {
	case humidity < 13 && tempF <= 112:
		return base - (13-float64(humidity))/4*math.Sqrt((17-math.Abs(float64(tempF)-95))/17)
	case humidity > 85 && tempF <= 87:
		return base + (float64(humidity)-85)/10*((87-float64(tempF))/5)
	default:
		return base
	}
}

// FeelsLikeF computes the feels-like temperature in Fahrenheit.
func FeelsLikeF[Temp, Humidity, WindSpeed constraints.Number](tempF Temp, humidity Humidity, windSpeedMPH WindSpeed) float64 {
	switch {
	case tempF <= 50 && windSpeedMPH > 3:
		return WindChillF(tempF, windSpeedMPH)
	case tempF > 68:
		return HeatIndexF(tempF, humidity)
	default:
		return float64(tempF)
	}
}
