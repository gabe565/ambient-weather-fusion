package ambientweather

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func (s *Server) ConnectMQTT(ctx context.Context) error {
	log := slog.With("url", s.conf.MQTTURL.String())
	cliCfg := autopaho.ClientConfig{
		ServerUrls:            []*url.URL{s.conf.MQTTURL.URL},
		TlsCfg:                &tls.Config{InsecureSkipVerify: s.conf.MQTTInsecureSkipVerify}, //nolint:gosec
		KeepAlive:             s.conf.MQTTKeepAlive,
		SessionExpiryInterval: s.conf.MQTTSessionExpiry,
		OnConnectionUp: func(_ *autopaho.ConnectionManager, _ *paho.Connack) {
			log.Info("Connected to MQTT")
			if err := s.PublishStatus(ctx, true); err != nil {
				log.Error("Failed to publish status message", "error", err)
			}
		},
		OnConnectError: func(err error) {
			log.Error("Failed to connect to MQTT", "error", err)
		},
		ConnectUsername: s.conf.MQTTUsername,
		ConnectPassword: []byte(s.conf.MQTTPassword),
		WillMessage: &paho.WillMessage{
			QoS:     1,
			Retain:  true,
			Topic:   s.conf.TopicPrefix + "/status",
			Payload: []byte("offline"),
		},
		ClientConfig: paho.ClientConfig{
			ClientID: s.conf.TopicPrefix,
			OnClientError: func(err error) {
				log.Error("Client error", "error", err)
			},
			OnServerDisconnect: func(d *paho.Disconnect) {
				var disconnectLog *slog.Logger
				if d.Properties != nil {
					disconnectLog = log.With("reason", d.Properties.ReasonString)
				} else {
					disconnectLog = log.With("reason", d.ReasonCode)
				}
				disconnectLog.Info("Server requested disconnect")
			},
		},
	}

	var err error
	s.mqtt, err = autopaho.NewConnection(context.Background(), cliCfg)
	if err != nil {
		return err
	}

	if err = s.mqtt.AwaitConnection(ctx); err != nil {
		_ = s.mqtt.Disconnect(ctx)
		s.mqtt = nil
		return err
	}

	return nil
}
