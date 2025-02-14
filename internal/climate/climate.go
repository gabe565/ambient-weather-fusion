package climate

import "math"

func CtoF(c float64) float64 {
	return c*9/5 + 32
}

func FtoC(f float64) float64 {
	return (f - 32) * 5 / 9
}

func KPHtoMPH(kph float64) float64 {
	return kph * 0.621371192
}

const (
	magnusA = 17.27
	magnusB = 237.7
)

func DewPointC(tempC, humidity float64) float64 {
	switch {
	case tempC == 0, humidity == 0:
		return 0
	default:
		humidity = min(humidity, 100)
		g := magnusA*tempC/(magnusB+tempC) + math.Log(humidity/100.0)
		return magnusB * g / (magnusA - g)
	}
}

func DewPointF(tempF, humidity float64) float64 {
	switch {
	case tempF == 0, humidity == 0:
		return 0
	default:
		tempC := FtoC(tempF)
		dewPoint := DewPointC(tempC, humidity)
		return CtoF(dewPoint)
	}
}

func windChillF(tempF, windSpeedMPH float64) float64 {
	exp := math.Pow(windSpeedMPH, 0.16)
	return 35.74 + 0.6215*tempF - 35.75*exp + 0.4275*tempF*exp
}

func heatIndexF(tempF, humidity float64) float64 {
	result := 0.5 * (tempF + 61 + (tempF-68)*1.2 + humidity*0.094)
	if tempF < 80 {
		return result
	}

	base := -42.379 +
		2.04901523*tempF +
		10.14333127*humidity +
		-0.22475541*tempF*humidity +
		-0.00683783*tempF*tempF +
		-0.05481717*humidity*humidity +
		0.00122874*tempF*tempF*humidity +
		0.00085282*tempF*humidity*humidity +
		-0.00000199*tempF*tempF*humidity*humidity
	switch {
	case humidity < 13 && tempF <= 112:
		return base - (13-humidity)/4*math.Sqrt((17-math.Abs(tempF-95))/17)
	case humidity > 85 && tempF <= 87:
		return base + (humidity-85)/10*((87-tempF)/5)
	default:
		return base
	}
}

func FeelsLikeF(tempF, humidity, windSpeedMPH float64) float64 {
	switch {
	case tempF == 0, humidity == 0, windSpeedMPH == 0:
		return 0
	case tempF < 50 && windSpeedMPH > 3:
		return windChillF(tempF, windSpeedMPH)
	case tempF > 68:
		return heatIndexF(tempF, humidity)
	default:
		return tempF
	}
}

func FeelsLikeC(tempC, humidity, windSpeedKPH float64) float64 {
	switch {
	case tempC == 0, humidity == 0, windSpeedKPH == 0:
		return 0
	default:
		tempF := CtoF(tempC)
		windSpeedMPH := KPHtoMPH(windSpeedKPH)
		feelsLikeF := FeelsLikeF(tempF, humidity, windSpeedMPH)
		return FtoC(feelsLikeF)
	}
}
