package ambientweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"slices"
	"time"

	"gabe565.com/ambient-weather-fusion/internal/config"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrNoEntries       = errors.New("no entries passed sanitization")
)

func expectedTokens() []json.Token {
	return []json.Token{
		json.Delim('{'),
		"data",
		json.Delim('['),
	}
}

var ErrInvalidType = errors.New("invalid type")

func Process(ctx context.Context, conf *config.Config, client *autopaho.ConnectionManager) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, conf.BuildURL().String(), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	decoder := json.NewDecoder(res.Body)
	for _, expect := range expectedTokens() {
		got, err := decoder.Token()
		if err != nil {
			return err
		}
		if got != expect {
			return fmt.Errorf("%w: got %s, expected %s", ErrInvalidResponse, got, expect)
		}
	}

	entries := make([]Data, 0, conf.Limit)
	for decoder.More() {
		var entry Data
		if err := decoder.Decode(&entry); err != nil {
			continue
		}

		if entry.Info.Indoor || entry.LastData.TempF == 0 {
			continue
		}

		var t time.Time
		switch {
		case entry.LastData.CreatedAt != 0:
			t = time.UnixMilli(entry.LastData.CreatedAt)
		case entry.LastData.DateUTC != 0:
			t = time.UnixMilli(entry.LastData.DateUTC)
		default:
			continue
		}

		if time.Since(t) > conf.MaxReadingAge {
			continue
		}

		ComputeValues(&entry)
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		return ErrNoEntries
	}

	var group errgroup.Group
	group.SetLimit(4)

	publish := func(topic string, value any) {
		if value == nil {
			return
		}
		group.Go(func() error {
			var b []byte
			switch topic {
			case "last_rain":
				if i, ok := value.(*int64); ok {
					ts := time.UnixMilli(*i)
					b = []byte(ts.Format(time.RFC3339))
				} else {
					return fmt.Errorf("%w for last_rain: %d", ErrInvalidType, value)
				}
			default:
				if b, err = json.Marshal(value); err != nil {
					return err
				}
			}

			topic := path.Join(conf.TopicPrefix, topic)
			slog.Debug("Publishing data", "topic", topic, "value", string(b))
			_, err = client.Publish(ctx, &paho.Publish{
				QoS:     1,
				Retain:  true,
				Topic:   topic,
				Payload: b,
			})
			return err
		})
	}

	publish("temperature", computeMedian(entries, func(data Data) float64 { return data.LastData.TempF }))
	publish("humidity", computeMedian(entries, func(data Data) int { return data.LastData.Humidity }))
	publish("wind_speed", computeMedian(entries, func(data Data) float64 { return data.LastData.WindSpeedMPH }))
	publish("wind_gust", computeMedian(entries, func(data Data) float64 { return data.LastData.WindGustMPH }))
	publish("max_daily_gust", computeMedian(entries, func(data Data) float64 { return data.LastData.MaxDailyGust }))
	publish("uv_index", computeMedian(entries, func(data Data) int { return data.LastData.UV }))
	publish("solar_radiation", computeMedian(entries, func(data Data) float64 { return data.LastData.SolarRadiation }))
	publish("hourly_rain", computeMedian(entries, func(data Data) float64 { return data.LastData.HourlyRainIn }))
	publish("daily_rain", computeMedian(entries, func(data Data) float64 { return data.LastData.DailyRainIn }))
	publish("weekly_rain", computeMedian(entries, func(data Data) float64 { return data.LastData.WeeklyRainIn }))
	publish("monthly_rain", computeMedian(entries, func(data Data) float64 { return data.LastData.MonthlyRainIn }))
	publish("relative_pressure", computeMedian(entries, func(data Data) float64 { return data.LastData.PressureRelativeIn }))
	publish("absolute_pressure", computeMedian(entries, func(data Data) float64 { return data.LastData.PressureAbsoluteIn }))
	publish("last_rain", computeMedian(entries, func(data Data) int64 { return data.LastData.LastRain }))
	publish("feels_like", computeMedian(entries, func(data Data) float64 { return data.LastData.FeelsLike }))
	publish("dew_point", computeMedian(entries, func(data Data) float64 { return data.LastData.DewPoint }))
	return group.Wait()
}

func computeMedian[V int | int64 | float64](inputs []Data, fn func(Data) V) *V {
	vals := make([]V, 0, len(inputs))
	for _, entry := range inputs {
		val := fn(entry)
		vals = append(vals, val)
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
