package discovery

type Platform string

const (
	PlatformSensor Platform = "sensor"
)

type Unit string

const (
	UnitFahrenheit      Unit = "°F"
	UnitPercent         Unit = "%"
	UnitMPH             Unit = "mph"
	UnitIndex           Unit = "index"
	UnitInches          Unit = "in"
	UnitInHg            Unit = "inHg"
	UnitInchesPerHour   Unit = "in/h"
	UnitWattsPerSqMeter Unit = "W/m²"
)

type DeviceClass string

const (
	DeviceClassTemperature            DeviceClass = "temperature"
	DeviceClassHumidity               DeviceClass = "humidity"
	DeviceClassWindSpeed              DeviceClass = "wind_speed"
	DeviceClassPrecipitation          DeviceClass = "precipitation"
	DeviceClassPrecipitationIntensity DeviceClass = "precipitation_intensity"
	DeviceClassPressure               DeviceClass = "pressure"
	DeviceClassTimestamp              DeviceClass = "timestamp"
	DeviceClassIrradiance             DeviceClass = "irradiance"
)

type StateClass string

const (
	StateClassMeasurement StateClass = "measurement"
	StateClassTotal       StateClass = "total"
)

type Topic string

const (
	TopicTemperature      Topic = "temperature"
	TopicHumidity         Topic = "humidity"
	TopicWindSpeed        Topic = "wind_speed"
	TopicWindGust         Topic = "wind_gust"
	TopicMaxDailyGust     Topic = "max_daily_gust"
	TopicUVIndex          Topic = "uv_index"
	TopicSolarRadiation   Topic = "solar_radiation"
	TopicHourlyRain       Topic = "hourly_rain"
	TopicDailyRain        Topic = "daily_rain"
	TopicWeeklyRain       Topic = "weekly_rain"
	TopicMonthlyRain      Topic = "monthly_rain"
	TopicRelativePressure Topic = "relative_pressure"
	TopicAbsolutePressure Topic = "absolute_pressure"
	TopicLastRain         Topic = "last_rain"
	TopicFeelsLike        Topic = "feels_like"
	TopicDewPoint         Topic = "dew_point"
)
