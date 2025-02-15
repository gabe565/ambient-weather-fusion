package mqtt

import (
	"context"
	"encoding/json"
	"log/slog"
	"path"

	"gabe565.com/ambient-weather-fusion/internal/config"
	"gabe565.com/utils/cobrax"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/spf13/cobra"
)

func PublishDiscovery(ctx context.Context, cmd *cobra.Command, conf *config.Config, client *autopaho.ConnectionManager) error {
	payload := generateDiscoveryPayload(cmd, conf)

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	topic := DiscoveryTopic(conf)
	slog.Debug("Publishing discovery payload", "topic", topic)
	_, err = client.Publish(ctx, &paho.Publish{
		QoS:     1,
		Retain:  true,
		Topic:   topic,
		Payload: b,
	})
	return err
}

func DiscoveryTopic(conf *config.Config) string {
	return path.Join(conf.DiscoveryPrefix, "device", conf.TopicPrefix, "config")
}

const (
	name                      = "name"
	unitOfMeasurement         = "unit_of_measurement"
	deviceClass               = "device_class"
	stateClass                = "state_class"
	suggestedDisplayPrecision = "suggested_display_precision"
	enabledByDefault          = "enabled_by_default"

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

func generateDiscoveryPayload(cmd *cobra.Command, conf *config.Config) map[string]any {
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
		sensor["object_id"] = conf.TopicPrefix + "_" + topic
		sensor["unique_id"] = conf.TopicPrefix + "_" + topic
		sensor["state_topic"] = path.Join(conf.TopicPrefix, topic)
	}

	payload := map[string]any{
		"availability": []map[string]any{
			{"topic": conf.TopicPrefix + "/status"},
		},
		"device": map[string]any{
			"ids": []string{conf.TopicPrefix},
			name:  conf.DeviceName,
			"sw":  cobrax.GetVersion(cmd),
		},
		"origin": map[string]any{
			name:  "Ambient Fusion",
			"sw":  cobrax.GetVersion(cmd),
			"url": "https://github.com/gabe565/ambient-fusion",
		},
		"components": components,
	}

	return payload
}
