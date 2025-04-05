package discovery

import "gabe565.com/ambient-weather-fusion/internal/config"

const (
	availabilityTopic = "avty_t"
	device            = "dev"
	identifiers       = "ids"
	swVersion         = "sw"
	origin            = "o"
	supportURL        = "url"
	stateTopic        = "stat_t"
	components        = "cmps"

	name                      = "name"
	platform                  = "p"
	objectID                  = "obj_id"
	uniqueID                  = "uniq_id"
	valueTemplate             = "val_tpl"
	unitOfMeasurement         = "unit_of_meas"
	deviceClass               = "dev_cla"
	stateClass                = "stat_cla"
	suggestedDisplayPrecision = "sug_dsp_prc"
	enabledByDefault          = "en"
	icon                      = "ic"

	unitFahrenheit      = "°F"
	unitPercent         = "%"
	unitMPH             = "mph"
	unitInches          = "in"
	unitInHg            = "inHg"
	unitInchesPerHour   = "in/h"
	unitWattsPerSqMeter = "W/m²"

	deviceClassTemperature            = "temperature"
	deviceClassHumidity               = "humidity"
	deviceClassWindSpeed              = "wind_speed"
	deviceClassPrecipitation          = "precipitation"
	deviceClassPrecipitationIntensity = "precipitation_intensity"
	deviceClassPressure               = "pressure"
	deviceClassTimestamp              = "timestamp"
	deviceClassIrradiance             = "irradiance"

	stateClassMeasurement = "measurement"
	stateClassTotal       = "total"

	TopicTemperature      = "temperature"
	TopicHumidity         = "humidity"
	TopicWindSpeed        = "wind_speed"
	TopicWindGust         = "wind_gust"
	TopicMaxDailyGust     = "max_daily_gust"
	TopicUVIndex          = "uv_index"
	TopicSolarRadiation   = "solar_radiation"
	TopicHourlyRain       = "hourly_rain"
	TopicDailyRain        = "daily_rain"
	TopicWeeklyRain       = "weekly_rain"
	TopicMonthlyRain      = "monthly_rain"
	TopicRelativePressure = "relative_pressure"
	TopicAbsolutePressure = "absolute_pressure"
	TopicLastRain         = "last_rain"
	TopicFeelsLike        = "feels_like"
	TopicDewPoint         = "dew_point"
)

func NewPayload(conf *config.Config, version string) map[string]any { //nolint:funlen
	c := map[string]map[string]any{
		TopicTemperature: {
			unitOfMeasurement:         unitFahrenheit,
			deviceClass:               deviceClassTemperature,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		TopicHumidity: {
			unitOfMeasurement:         unitPercent,
			deviceClass:               deviceClassHumidity,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		TopicWindSpeed: {
			unitOfMeasurement:         unitMPH,
			deviceClass:               deviceClassWindSpeed,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		TopicWindGust: {
			name:                      "Wind gust",
			unitOfMeasurement:         unitMPH,
			deviceClass:               deviceClassWindSpeed,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		TopicMaxDailyGust: {
			name:                      "Max daily gust",
			unitOfMeasurement:         unitMPH,
			deviceClass:               deviceClassWindSpeed,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		TopicUVIndex: {
			name:                      "UV index",
			unitOfMeasurement:         "index",
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		TopicSolarRadiation: {
			unitOfMeasurement:         unitWattsPerSqMeter,
			deviceClass:               deviceClassIrradiance,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
			enabledByDefault:          false,
		},
		TopicHourlyRain: {
			name:                      "Hourly rain",
			unitOfMeasurement:         unitInchesPerHour,
			deviceClass:               deviceClassPrecipitationIntensity,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 2,
		},
		TopicDailyRain: {
			name:                      "Daily rain",
			unitOfMeasurement:         unitInches,
			deviceClass:               deviceClassPrecipitation,
			stateClass:                stateClassTotal,
			suggestedDisplayPrecision: 2,
		},
		TopicWeeklyRain: {
			name:                      "Weekly rain",
			unitOfMeasurement:         unitInches,
			deviceClass:               deviceClassPrecipitation,
			stateClass:                stateClassTotal,
			suggestedDisplayPrecision: 2,
			enabledByDefault:          false,
		},
		TopicMonthlyRain: {
			name:                      "Monthly rain",
			unitOfMeasurement:         unitInches,
			deviceClass:               deviceClassPrecipitation,
			stateClass:                stateClassTotal,
			suggestedDisplayPrecision: 2,
			enabledByDefault:          false,
		},
		TopicRelativePressure: {
			name:                      "Relative pressure",
			unitOfMeasurement:         unitInHg,
			deviceClass:               deviceClassPressure,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 2,
		},
		TopicAbsolutePressure: {
			name:                      "Absolute pressure",
			unitOfMeasurement:         unitInHg,
			deviceClass:               deviceClassPressure,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 2,
			enabledByDefault:          false,
		},
		TopicLastRain: {
			name:             "Last rain",
			deviceClass:      deviceClassTimestamp,
			enabledByDefault: false,
			icon:             "mdi:water",
		},
		TopicFeelsLike: {
			name:                      "Feels like",
			unitOfMeasurement:         unitFahrenheit,
			deviceClass:               deviceClassTemperature,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		TopicDewPoint: {
			name:                      "Dew point",
			unitOfMeasurement:         unitFahrenheit,
			deviceClass:               deviceClassTemperature,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
	}

	for topic, sensor := range c {
		sensor[platform] = "sensor"
		sensor[objectID] = conf.BaseTopic + "_" + topic
		sensor[uniqueID] = conf.BaseTopic + "_" + topic
		sensor[valueTemplate] = "{{ value_json." + topic + " }}"
	}

	payload := map[string]any{
		availabilityTopic: conf.BaseTopic + "/status",
		device: map[string]any{
			identifiers: conf.BaseTopic,
			name:        conf.HADeviceName,
			swVersion:   version,
		},
		origin: map[string]any{
			name:       "Ambient Weather Fusion",
			swVersion:  version,
			supportURL: "https://github.com/gabe565/ambient-weather-fusion",
		},
		stateTopic: conf.BaseTopic,
		components: c,
	}

	return payload
}
