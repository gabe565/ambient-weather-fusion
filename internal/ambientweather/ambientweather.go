package ambientweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"gabe565.com/ambient-weather-fusion/internal/config"
	"gabe565.com/ambient-weather-fusion/internal/mqtt"
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

	for topic, value := range generatePayloads(entries) {
		if value == nil {
			continue
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

			slog.Debug("Publishing data", "topic", topic, "value", string(b))
			_, err = client.Publish(ctx, &paho.Publish{
				QoS:     1,
				Retain:  true,
				Topic:   mqtt.DataTopic(conf, topic),
				Payload: b,
			})
			return err
		})
	}

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

func generatePayloads(entries []Data) map[string]any {
	return map[string]any{
		"temperature":       computeMedian(entries, func(data Data) float64 { return data.LastData.TempF }),
		"humidity":          computeMedian(entries, func(data Data) int { return data.LastData.Humidity }),
		"wind_speed":        computeMedian(entries, func(data Data) float64 { return data.LastData.WindSpeedMPH }),
		"wind_gust":         computeMedian(entries, func(data Data) float64 { return data.LastData.WindGustMPH }),
		"max_daily_gust":    computeMedian(entries, func(data Data) float64 { return data.LastData.MaxDailyGust }),
		"uv_index":          computeMedian(entries, func(data Data) int { return data.LastData.UV }),
		"solar_radiation":   computeMedian(entries, func(data Data) float64 { return data.LastData.SolarRadiation }),
		"hourly_rain":       computeMedian(entries, func(data Data) float64 { return data.LastData.HourlyRainIn }),
		"daily_rain":        computeMedian(entries, func(data Data) float64 { return data.LastData.DailyRainIn }),
		"weekly_rain":       computeMedian(entries, func(data Data) float64 { return data.LastData.WeeklyRainIn }),
		"monthly_rain":      computeMedian(entries, func(data Data) float64 { return data.LastData.MonthlyRainIn }),
		"relative_pressure": computeMedian(entries, func(data Data) float64 { return data.LastData.PressureRelativeIn }),
		"absolute_pressure": computeMedian(entries, func(data Data) float64 { return data.LastData.PressureAbsoluteIn }),
		"last_rain":         computeMedian(entries, func(data Data) int64 { return data.LastData.LastRain }),
		"feels_like":        computeMedian(entries, func(data Data) float64 { return data.LastData.FeelsLike }),
		"dew_point":         computeMedian(entries, func(data Data) float64 { return data.LastData.DewPoint }),
	}
}
