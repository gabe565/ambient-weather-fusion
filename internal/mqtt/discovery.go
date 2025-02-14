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
	availability := []map[string]any{{"topic": conf.TopicPrefix + "/status"}}
	origin := map[string]any{
		"name": "Ambient Fusion",
		"sw":   cobrax.GetVersion(cmd),
		"url":  "https://github.com/gabe565/ambient-fusion",
	}
	device := map[string]any{
		"identifiers": []string{conf.TopicPrefix},
		"name":        conf.DeviceName,
		"sw_version":  cobrax.GetVersion(cmd),
	}

	var group errgroup.Group
	group.SetLimit(4)

	publish := func(topic string, data map[string]any) {
		group.Go(func() error {
			data["availability"] = availability
			data["origin"] = origin
			data["device"] = device
			data["object_id"] = conf.TopicPrefix + "_" + topic
			data["unique_id"] = conf.TopicPrefix + "_" + topic
			data["state_topic"] = path.Join(conf.TopicPrefix, topic)

			b, err := json.Marshal(data)
			if err != nil {
				return err
			}

			topic := path.Join(conf.DiscoveryPrefix, "sensor", conf.TopicPrefix, topic, "config")
			slog.Debug("Publishing discovery payload", "topic", topic)
			_, err = client.Publish(ctx, &paho.Publish{
				QoS:     1,
				Retain:  true,
				Topic:   topic,
				Payload: b,
			})
			return err
		})
	}

	publish("temperature", map[string]any{
		"unit_of_measurement":         "°F",
		"device_class":                "temperature",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})
	publish("humidity", map[string]any{
		"unit_of_measurement":         "%",
		"device_class":                "humidity",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})
	publish("wind_speed", map[string]any{
		"unit_of_measurement":         "mph",
		"device_class":                "wind_speed",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})
	publish("wind_gust", map[string]any{
		"name":                        "Wind gust",
		"unit_of_measurement":         "mph",
		"device_class":                "wind_speed",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})
	publish("max_daily_gust", map[string]any{
		"name":                        "Max daily gust",
		"unit_of_measurement":         "mph",
		"device_class":                "wind_speed",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})
	publish("uv_index", map[string]any{
		"name":                        "UV index",
		"unit_of_measurement":         "index",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})
	publish("solar_radiation", map[string]any{
		"unit_of_measurement":         "W/m²",
		"device_class":                "irradiance",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
		"enabled_by_default":          false,
	})
	publish("hourly_rain", map[string]any{
		"name":                        "Hourly rain",
		"unit_of_measurement":         "in/h",
		"device_class":                "precipitation_intensity",
		"state_class":                 "measurement",
		"suggested_display_precision": 2,
	})
	publish("daily_rain", map[string]any{
		"name":                        "Daily rain",
		"unit_of_measurement":         "in",
		"device_class":                "precipitation",
		"state_class":                 "total",
		"suggested_display_precision": 2,
	})
	publish("weekly_rain", map[string]any{
		"name":                        "Weekly rain",
		"unit_of_measurement":         "in",
		"device_class":                "precipitation",
		"state_class":                 "total",
		"suggested_display_precision": 2,
		"enabled_by_default":          false,
	})
	publish("monthly_rain", map[string]any{
		"name":                        "Monthly rain",
		"unit_of_measurement":         "in",
		"device_class":                "precipitation",
		"state_class":                 "total",
		"suggested_display_precision": 2,
		"enabled_by_default":          false,
	})
	publish("relative_pressure", map[string]any{
		"name":                        "Relative pressure",
		"unit_of_measurement":         "inHg",
		"device_class":                "pressure",
		"state_class":                 "measurement",
		"suggested_display_precision": 2,
	})
	publish("absolute_pressure", map[string]any{
		"name":                        "Absolute pressure",
		"unit_of_measurement":         "inHg",
		"device_class":                "pressure",
		"state_class":                 "measurement",
		"suggested_display_precision": 2,
		"enabled_by_default":          false,
	})
	publish("last_rain", map[string]any{
		"name":               "Last rain",
		"device_class":       "timestamp",
		"enabled_by_default": false,
	})
	publish("feels_like", map[string]any{
		"name":                        "Feels like",
		"unit_of_measurement":         "°F",
		"device_class":                "temperature",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})
	publish("dew_point", map[string]any{
		"name":                        "Dew point",
		"unit_of_measurement":         "°F",
		"device_class":                "temperature",
		"state_class":                 "measurement",
		"suggested_display_precision": 1,
	})

	return group.Wait()
}
