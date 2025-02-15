package ambientweather

import (
	"gabe565.com/ambient-weather-fusion/internal/climate"
)

func ComputeValues(data *Data) {
	if data.LastData.TempF == nil || data.LastData.Humidity == nil {
		return
	}

	if data.LastData.DewPoint == nil {
		dewPoint := climate.DewPointF(
			*data.LastData.TempF,
			float64(*data.LastData.Humidity),
		)
		data.LastData.DewPoint = &dewPoint
	}
	if data.LastData.FeelsLike == nil && data.LastData.WindGustMPH != nil {
		feelsLike := climate.FeelsLikeF(
			*data.LastData.TempF,
			float64(*data.LastData.Humidity),
			*data.LastData.WindGustMPH,
		)
		data.LastData.FeelsLike = &feelsLike
	}
}
