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
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/httpx"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/spf13/cobra"
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

func Process(ctx context.Context, cmd *cobra.Command, conf *config.Config, client *autopaho.ConnectionManager) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, conf.BuildURL().String(), nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Transport: httpx.NewUserAgentTransport(nil, cobrax.BuildUserAgent(cmd)),
		Timeout:   time.Minute,
	}

	res, err := httpClient.Do(req)
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

		if entry.Info.Indoor == nil || *entry.Info.Indoor || entry.LastData.TempF == nil {
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
			case mqtt.TopicLastRain:
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

func computeMedian[V int | int64 | float64](inputs []Data, fn func(Data) *V) *V {
	vals := make([]V, 0, len(inputs))
	for _, entry := range inputs {
		if val := fn(entry); val != nil {
			vals = append(vals, *val)
		}
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
		mqtt.TopicTemperature:      computeMedian(entries, func(data Data) *float64 { return data.LastData.TempF }),
		mqtt.TopicHumidity:         computeMedian(entries, func(data Data) *int { return data.LastData.Humidity }),
		mqtt.TopicWindSpeed:        computeMedian(entries, func(data Data) *float64 { return data.LastData.WindSpeedMPH }),
		mqtt.TopicWindGust:         computeMedian(entries, func(data Data) *float64 { return data.LastData.WindGustMPH }),
		mqtt.TopicMaxDailyGust:     computeMedian(entries, func(data Data) *float64 { return data.LastData.MaxDailyGust }),
		mqtt.TopicUVIndex:          computeMedian(entries, func(data Data) *int { return data.LastData.UV }),
		mqtt.TopicSolarRadiation:   computeMedian(entries, func(data Data) *float64 { return data.LastData.SolarRadiation }),
		mqtt.TopicHourlyRain:       computeMedian(entries, func(data Data) *float64 { return data.LastData.HourlyRainIn }),
		mqtt.TopicDailyRain:        computeMedian(entries, func(data Data) *float64 { return data.LastData.DailyRainIn }),
		mqtt.TopicWeeklyRain:       computeMedian(entries, func(data Data) *float64 { return data.LastData.WeeklyRainIn }),
		mqtt.TopicMonthlyRain:      computeMedian(entries, func(data Data) *float64 { return data.LastData.MonthlyRainIn }),
		mqtt.TopicRelativePressure: computeMedian(entries, func(data Data) *float64 { return data.LastData.PressureRelativeIn }),
		mqtt.TopicAbsolutePressure: computeMedian(entries, func(data Data) *float64 { return data.LastData.PressureAbsoluteIn }),
		mqtt.TopicLastRain:         computeMedian(entries, func(data Data) *int64 { return data.LastData.LastRain }),
		mqtt.TopicFeelsLike:        computeMedian(entries, func(data Data) *float64 { return data.LastData.FeelsLike }),
		mqtt.TopicDewPoint:         computeMedian(entries, func(data Data) *float64 { return data.LastData.DewPoint }),
	}
}

func Cleanup(ctx context.Context, conf *config.Config, client *autopaho.ConnectionManager) error {
	var group errgroup.Group
	group.SetLimit(4)
	for topic := range generatePayloads(nil) {
		group.Go(func() error {
			_, err := client.Publish(ctx, &paho.Publish{
				Topic: mqtt.DataTopic(conf, topic),
			})
			return err
		})
	}
	return group.Wait()
}
