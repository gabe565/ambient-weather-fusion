package ambientweather

import (
	"gabe565.com/ambient-weather-fusion/internal/climate"
)

func ComputeValues(data *Data) {
	if data.LastData.DewPoint == 0 {
		data.LastData.DewPoint = climate.DewPointF(
			data.LastData.TempF,
			float64(data.LastData.Humidity),
		)
	}
	if data.LastData.FeelsLike == 0 {
		data.LastData.FeelsLike = climate.FeelsLikeF(
			data.LastData.TempF,
			float64(data.LastData.Humidity),
			data.LastData.WindSpeedMPH,
		)
	}
}
