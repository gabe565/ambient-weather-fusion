package ambientweather

import (
	"context"
	"encoding/json"
	"log/slog"
	"path"

	"gabe565.com/ambient-weather-fusion/internal/ambientweather/discovery"
	"github.com/eclipse/paho.golang/paho"
)

func (s *Server) PublishDiscovery(ctx context.Context) error {
	b, err := json.Marshal(discovery.NewPayload(s.conf, s.version))
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
