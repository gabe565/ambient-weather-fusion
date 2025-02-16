package ambientweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gabe565.com/ambient-weather-fusion/internal/config"
	"gabe565.com/ambient-weather-fusion/internal/location"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func NewServer(conf *config.Config, options ...Option) *Server {
	s := &Server{
		conf: conf,
		http: &http.Client{Timeout: time.Minute},
	}
	for _, option := range options {
		option(s)
	}
	return s
}

type Server struct {
	conf      *config.Config
	mqtt      *autopaho.ConnectionManager
	http      *http.Client
	version   string
	userAgent string
}

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

func (s *Server) BuildURL() *url.URL {
	u := *s.conf.RequestURL.URL
	q := u.Query()
	lat1, lon1 := location.Shift(s.conf.Latitude, s.conf.Longitude, -s.conf.Radius, -s.conf.Radius)
	lat2, lon2 := location.Shift(s.conf.Latitude, s.conf.Longitude, s.conf.Radius, s.conf.Radius)
	q.Set("$publicBox[0][0]", strconv.FormatFloat(lon1, 'f', -1, 64))
	q.Set("$publicBox[0][1]", strconv.FormatFloat(lat1, 'f', -1, 64))
	q.Set("$publicBox[1][0]", strconv.FormatFloat(lon2, 'f', -1, 64))
	q.Set("$publicBox[1][1]", strconv.FormatFloat(lat2, 'f', -1, 64))
	q.Set("$limit", strconv.Itoa(s.conf.Limit))
	u.RawQuery = q.Encode()
	return &u
}

func (s *Server) FetchData(ctx context.Context) ([]Data, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.BuildURL().String(), nil)
	if err != nil {
		return nil, err
	}
	if s.userAgent != "" {
		req.Header.Set("User-Agent", s.userAgent)
	}

	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", ErrUpstream, res.Status)
	}

	decoder := json.NewDecoder(res.Body)
	for _, expect := range expectedTokens() {
		got, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if got != expect {
			return nil, fmt.Errorf("%w: got %s, expected %s", ErrInvalidResponse, got, expect)
		}
	}

	data := make([]Data, 0, s.conf.Limit)
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

		if time.Since(t) > s.conf.MaxReadingAge {
			continue
		}

		data = append(data, entry)
	}

	if len(data) == 0 {
		return nil, ErrNoEntries
	}
	return data, nil
}

func (s *Server) PublishStatus(ctx context.Context, online bool) error {
	payload := "online"
	if !online {
		payload = "offline"
	}

	slog.Debug("Publishing status payload", "topic", s.conf.TopicPrefix, "payload", payload)
	_, err := s.mqtt.Publish(ctx, &paho.Publish{
		QoS:     1,
		Retain:  true,
		Topic:   s.conf.TopicPrefix + "/status",
		Payload: []byte(payload),
	})
	return err
}

func (s *Server) PublishData(ctx context.Context, retain bool, payload *Payload) error {
	var b []byte
	if payload != nil {
		var err error
		if b, err = json.Marshal(payload); err != nil {
			return err
		}
	}

	slog.Debug("Publishing data payload", "topic", s.conf.TopicPrefix, "retain", retain, "payload", string(b))
	_, err := s.mqtt.Publish(ctx, &paho.Publish{
		QoS:     1,
		Retain:  retain,
		Topic:   s.conf.TopicPrefix,
		Payload: b,
	})
	return err
}

func (s *Server) Close(ctx context.Context) error {
	if s.mqtt == nil {
		return nil
	}

	defer func() {
		s.mqtt = nil
	}()

	return errors.Join(
		s.PublishStatus(ctx, false),
		s.PublishData(ctx, false, nil),
		s.mqtt.Disconnect(ctx),
	)
}

func (s *Server) Tick(ctx context.Context) error {
	data, err := s.FetchData(ctx)
	if err != nil {
		return err
	}

	return s.PublishData(ctx, true, NewPayload(data))
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.ConnectMQTT(ctx); err != nil {
		return err
	}

	if err := s.PublishDiscovery(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			ticker.Reset(5 * time.Minute)
			if err := s.Tick(ctx); err != nil {
				slog.Error("Failed to process ambient-weather data", "error", err)
			}
		}
	}
}
