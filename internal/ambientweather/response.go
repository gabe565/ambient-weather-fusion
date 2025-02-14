package ambientweather

type Response struct {
	Data []Data `json:"data"`
}

type Data struct {
	LastData LastData `json:"lastData"`
	Info     Info     `json:"info"`
}

type LastData struct {
	DateUTC            int64   `json:"dateutc"`
	TempF              float64 `json:"tempf"`
	Humidity           int     `json:"humidity"`
	WindSpeedMPH       float64 `json:"windspeedmph"`
	WindGustMPH        float64 `json:"windgustmph"`
	MaxDailyGust       float64 `json:"maxdailygust"`
	UV                 int     `json:"uv"`
	SolarRadiation     float64 `json:"solarradiation"`
	HourlyRainIn       float64 `json:"hourlyrainin"`
	DailyRainIn        float64 `json:"dailyrainin"`
	WeeklyRainIn       float64 `json:"weeklyrainin"`
	MonthlyRainIn      float64 `json:"monthlyrainin"`
	PressureRelativeIn float64 `json:"baromrelin"`
	PressureAbsoluteIn float64 `json:"baromabsin"`
	CreatedAt          int64   `json:"created_at"`
	LastRain           int64   `json:"lastRain"`
	FeelsLike          float64 `json:"feelsLike,omitempty"`
	DewPoint           float64 `json:"dewPoint,omitempty"`
}

type Info struct {
	Name   string `json:"name"`
	Indoor bool   `json:"indoor"`
	Slug   string `json:"slug"`
}
