package ambientweather

import "gabe565.com/ambient-weather-fusion/internal/climate"

type Response struct {
	Data []Data `json:"data"`
}

type Data struct {
	LastData LastData `json:"lastData"`
	Info     Info     `json:"info"`
}

type LastData struct {
	DateUTC            int64    `json:"dateutc"`
	TempF              *float64 `json:"tempf"`
	Humidity           *int     `json:"humidity"`
	WindSpeedMPH       *float64 `json:"windspeedmph"`
	WindGustMPH        *float64 `json:"windgustmph"`
	MaxDailyGust       *float64 `json:"maxdailygust"`
	UV                 *int     `json:"uv"`
	SolarRadiation     *float64 `json:"solarradiation"`
	HourlyRainIn       *float64 `json:"hourlyrainin"`
	DailyRainIn        *float64 `json:"dailyrainin"`
	WeeklyRainIn       *float64 `json:"weeklyrainin"`
	MonthlyRainIn      *float64 `json:"monthlyrainin"`
	PressureRelativeIn *float64 `json:"baromrelin"`
	PressureAbsoluteIn *float64 `json:"baromabsin"`
	CreatedAt          int64    `json:"created_at"`
	LastRain           *int64   `json:"lastRain"`
	FeelsLike          *float64 `json:"feelsLike,omitempty"`
	DewPoint           *float64 `json:"dewPoint,omitempty"`
}

func (l *LastData) GetFeelsLike() *float64 {
	if l.FeelsLike == nil && l.TempF != nil && l.Humidity != nil && l.WindGustMPH != nil {
		feelsLike := climate.FeelsLikeF(
			*l.TempF,
			float64(*l.Humidity),
			*l.WindGustMPH,
		)
		l.FeelsLike = &feelsLike
	}
	return l.FeelsLike
}

func (l *LastData) GetDewPoint() *float64 {
	if l.DewPoint == nil && l.TempF != nil && l.Humidity != nil {
		dewPoint := climate.DewPointF(
			*l.TempF,
			float64(*l.Humidity),
		)
		l.DewPoint = &dewPoint
	}
	return l.DewPoint
}

type Info struct {
	Name   string `json:"name"`
	Indoor *bool  `json:"indoor"`
	Slug   string `json:"slug"`
}
