package ambientweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"gabe565.com/ambient-weather-fusion/internal/config"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/httpx"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/spf13/cobra"
)

var (
	ErrUpstream        = errors.New("upstream returned an error")
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

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", ErrUpstream, res.Status)
	}

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

		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		return ErrNoEntries
	}

	b, err := json.Marshal(NewPayload(entries))
	if err != nil {
		return err
	}

	slog.Debug("Publishing data", "topic", conf.TopicPrefix, "value", string(b))
	_, err = client.Publish(ctx, &paho.Publish{
		QoS:     1,
		Retain:  true,
		Topic:   conf.TopicPrefix,
		Payload: b,
	})
	return err
}

func Cleanup(ctx context.Context, conf *config.Config, client *autopaho.ConnectionManager) error {
	_, err := client.Publish(ctx, &paho.Publish{
		Topic: conf.TopicPrefix,
	})
	return err
}
