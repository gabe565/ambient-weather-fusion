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
	"golang.org/x/sync/errgroup"
)

func PublishDiscovery(ctx context.Context, cmd *cobra.Command, conf *config.Config, client *autopaho.ConnectionManager) error {
	var group errgroup.Group
	group.SetLimit(4)

	for topic, data := range generateDiscoveryPayloads(cmd, conf) {
		group.Go(func() error {
			b, err := json.Marshal(data)
			if err != nil {
				return err
			}

			slog.Debug("Publishing discovery payload", "topic", topic)
			_, err = client.Publish(ctx, &paho.Publish{
				QoS:     1,
				Retain:  true,
				Topic:   DiscoveryTopic(conf, topic),
				Payload: b,
			})
			return err
		})
	}

	return group.Wait()
}

func DiscoveryTopic(conf *config.Config, topic string) string {
	return path.Join(conf.DiscoveryPrefix, "sensor", conf.TopicPrefix, topic, "config")
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
)

func generateDiscoveryPayloads(cmd *cobra.Command, conf *config.Config) map[string]map[string]any {
	payloads := map[string]map[string]any{
		"temperature": {
			unitOfMeasurement:         unitFahrenheit,
			deviceClass:               deviceClassTemperature,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		"humidity": {
			unitOfMeasurement:         unitPercent,
			deviceClass:               deviceClassHumidity,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		"wind_speed": {
			unitOfMeasurement:         unitMPH,
			deviceClass:               deviceClassWindSpeed,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		"wind_gust": {
			name:                      "Wind gust",
			unitOfMeasurement:         unitMPH,
			deviceClass:               deviceClassWindSpeed,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		"max_daily_gust": {
			name:                      "Max daily gust",
			unitOfMeasurement:         unitMPH,
			deviceClass:               deviceClassWindSpeed,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		"uv_index": {
			name:                      "UV index",
			unitOfMeasurement:         "index",
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		"solar_radiation": {
			unitOfMeasurement:         unitWattsPerSqMeter,
			deviceClass:               deviceClassIrradiance,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
			enabledByDefault:          false,
		},
		"hourly_rain": {
			name:                      "Hourly rain",
			unitOfMeasurement:         unitInchesPerHour,
			deviceClass:               deviceClassPrecipitationIntensity,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 2,
		},
		"daily_rain": {
			name:                      "Daily rain",
			unitOfMeasurement:         unitInches,
			deviceClass:               deviceClassPrecipitation,
			stateClass:                stateClassTotal,
			suggestedDisplayPrecision: 2,
		},
		"weekly_rain": {
			name:                      "Weekly rain",
			unitOfMeasurement:         unitInches,
			deviceClass:               deviceClassPrecipitation,
			stateClass:                stateClassTotal,
			suggestedDisplayPrecision: 2,
			enabledByDefault:          false,
		},
		"monthly_rain": {
			name:                      "Monthly rain",
			unitOfMeasurement:         unitInches,
			deviceClass:               deviceClassPrecipitation,
			stateClass:                stateClassTotal,
			suggestedDisplayPrecision: 2,
			enabledByDefault:          false,
		},
		"relative_pressure": {
			name:                      "Relative pressure",
			unitOfMeasurement:         unitInHg,
			deviceClass:               deviceClassPressure,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 2,
		},
		"absolute_pressure": {
			name:                      "Absolute pressure",
			unitOfMeasurement:         unitInHg,
			deviceClass:               deviceClassPressure,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 2,
			enabledByDefault:          false,
		},
		"last_rain": {
			name:             "Last rain",
			deviceClass:      deviceClassTimestamp,
			enabledByDefault: false,
		},
		"feels_like": {
			name:                      "Feels like",
			unitOfMeasurement:         unitFahrenheit,
			deviceClass:               deviceClassTemperature,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
		"dew_point": {
			name:                      "Dew point",
			unitOfMeasurement:         unitFahrenheit,
			deviceClass:               deviceClassTemperature,
			stateClass:                stateClassMeasurement,
			suggestedDisplayPrecision: 1,
		},
	}

	availability := []map[string]any{{
		"topic": conf.TopicPrefix + "/status",
	}}
	origin := map[string]any{
		name:  "Ambient Fusion",
		"sw":  cobrax.GetVersion(cmd),
		"url": "https://github.com/gabe565/ambient-fusion",
	}
	device := map[string]any{
		"identifiers": []string{conf.TopicPrefix},
		name:          conf.DeviceName,
		"sw_version":  cobrax.GetVersion(cmd),
	}

	for topic, sensor := range payloads {
		sensor["availability"] = availability
		sensor["origin"] = origin
		sensor["device"] = device
		sensor["object_id"] = conf.TopicPrefix + "_" + topic
		sensor["unique_id"] = conf.TopicPrefix + "_" + topic
		sensor["state_topic"] = path.Join(conf.TopicPrefix, topic)
	}
	return payloads
}
