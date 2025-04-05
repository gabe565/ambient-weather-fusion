package ambientweather

import (
	"context"
	"encoding/json"
	"log/slog"
	"path"

	"github.com/eclipse/paho.golang/paho"
)

func (s *Server) PublishDiscovery(ctx context.Context) error {
	b, err := json.Marshal(s.DiscoveryPayload())
	if err != nil {
		return err
	}

	topic := s.DiscoveryTopic()
	slog.Debug("Publishing discovery payload", "topic", topic)
	_, err = s.mqtt.Publish(ctx, &paho.Publish{
		QoS:     1,
		Retain:  true,
		Topic:   topic,
		Payload: b,
	})
	return err
}

func (s *Server) DiscoveryTopic() string {
	return path.Join(s.conf.HADiscoveryTopic, "device", s.conf.BaseTopic, "config")
}

const (
	name                      = "name"
	unitOfMeasurement         = "unit_of_measurement"
	deviceClass               = "device_class"
	stateClass                = "state_class"
	suggestedDisplayPrecision = "suggested_display_precision"
	enabledByDefault          = "enabled_by_default"
	icon                      = "icon"

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

func (s *Server) DiscoveryPayload() map[string]any { //nolint:funlen
	components := map[string]map[string]any{
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

	for topic, sensor := range components {
		sensor["platform"] = "sensor"
		sensor["object_id"] = s.conf.BaseTopic + "_" + topic
		sensor["unique_id"] = s.conf.BaseTopic + "_" + topic
		sensor["value_template"] = "{{ value_json." + topic + " }}"
	}

	payload := map[string]any{
		"availability": []map[string]any{
			{"topic": s.conf.BaseTopic + "/status"},
		},
		"device": map[string]any{
			"ids": []string{s.conf.BaseTopic},
			name:  s.conf.HADeviceName,
			"sw":  s.version,
		},
		"origin": map[string]any{
			name:  "Ambient Weather Fusion",
			"sw":  s.version,
			"url": "https://github.com/gabe565/ambient-weather-fusion",
		},
		"state_topic": s.conf.BaseTopic,
		"components":  components,
	}

	return payload
}
