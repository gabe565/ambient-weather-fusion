package ambientweather

import (
	"slices"
	"time"
)

type Payload struct {
	Temperature      *float64 `json:"temperature,omitempty"`
	Humidity         *int     `json:"humidity,omitempty"`
	WindSpeed        *float64 `json:"wind_speed,omitempty"`
	WindGust         *float64 `json:"wind_gust,omitempty"`
	MaxDailyGust     *float64 `json:"max_daily_gust,omitempty"`
	UVIndex          *int     `json:"uv_index,omitempty"`
	SolarRadiation   *float64 `json:"solar_radiation,omitempty"`
	HourlyRain       *float64 `json:"hourly_rain,omitempty"`
	DailyRain        *float64 `json:"daily_rain,omitempty"`
	WeeklyRain       *float64 `json:"weekly_rain,omitempty"`
	MonthlyRain      *float64 `json:"monthly_rain,omitempty"`
	RelativePressure *float64 `json:"relative_pressure,omitempty"`
	AbsolutePressure *float64 `json:"absolute_pressure,omitempty"`
	LastRain         *string  `json:"last_rain,omitempty"`
	FeelsLike        *float64 `json:"feels_like,omitempty"`
	DewPoint         *float64 `json:"dew_point,omitempty"`
}

func computeMedian[V int | int64 | float64](inputs []Data, fn func(Data) *V) *V {
	vals := make([]V, 0, len(inputs))
	for _, entry := range inputs {
		if val := fn(entry); val != nil {
			vals = append(vals, *val)
		}
	}

	slices.Sort(vals)

	var val V
	switch {
	case len(vals) == 0:
		return nil
	case len(vals) == 1:
		val = vals[0]
	case len(vals)%2 != 0:
		val = vals[len(vals)/2]
	default:
		val = (vals[len(vals)/2-1] + vals[len(vals)/2]) / 2
	}
	return &val
}

func NewPayload(entries []Data) *Payload {
	p := &Payload{
		Temperature:      computeMedian(entries, func(data Data) *float64 { return data.LastData.TempF }),
		Humidity:         computeMedian(entries, func(data Data) *int { return data.LastData.Humidity }),
		WindSpeed:        computeMedian(entries, func(data Data) *float64 { return data.LastData.WindSpeedMPH }),
		WindGust:         computeMedian(entries, func(data Data) *float64 { return data.LastData.WindGustMPH }),
		MaxDailyGust:     computeMedian(entries, func(data Data) *float64 { return data.LastData.MaxDailyGust }),
		UVIndex:          computeMedian(entries, func(data Data) *int { return data.LastData.UV }),
		SolarRadiation:   computeMedian(entries, func(data Data) *float64 { return data.LastData.SolarRadiation }),
		HourlyRain:       computeMedian(entries, func(data Data) *float64 { return data.LastData.HourlyRainIn }),
		DailyRain:        computeMedian(entries, func(data Data) *float64 { return data.LastData.DailyRainIn }),
		WeeklyRain:       computeMedian(entries, func(data Data) *float64 { return data.LastData.WeeklyRainIn }),
		MonthlyRain:      computeMedian(entries, func(data Data) *float64 { return data.LastData.MonthlyRainIn }),
		RelativePressure: computeMedian(entries, func(data Data) *float64 { return data.LastData.PressureRelativeIn }),
		AbsolutePressure: computeMedian(entries, func(data Data) *float64 { return data.LastData.PressureAbsoluteIn }),
		FeelsLike:        computeMedian(entries, func(data Data) *float64 { return data.LastData.GetFeelsLike() }),
		DewPoint:         computeMedian(entries, func(data Data) *float64 { return data.LastData.GetDewPoint() }),
	}

	if unix := computeMedian(entries, func(data Data) *int64 { return data.LastData.LastRain }); unix != nil {
		timestamp := time.UnixMilli(*unix).Format(time.RFC3339)
		p.LastRain = &timestamp
	}

	return p
}
