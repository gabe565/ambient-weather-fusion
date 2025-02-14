package mqtt

import (
	"context"
	"errors"
	"log/slog"
	"net/url"

	"gabe565.com/ambient-weather-fusion/internal/config"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func Connect(ctx context.Context, conf *config.Config) (*autopaho.ConnectionManager, error) {
	cliCfg := autopaho.ClientConfig{
		ServerUrls:            []*url.URL{conf.MQTTURL.URL},
		KeepAlive:             20,
		SessionExpiryInterval: 60,
		OnConnectionUp: func(client *autopaho.ConnectionManager, _ *paho.Connack) {
			slog.Info("Connected to MQTT")
			if _, err := client.Publish(ctx, &paho.Publish{
				QoS:     1,
				Retain:  true,
				Topic:   conf.TopicPrefix + "/status",
				Payload: []byte("online"),
			}); err != nil {
				slog.Error("Failed to publish status message", "error", err)
			}
		},
		OnConnectError: func(err error) {
			slog.Error("Failed to connect to MQTT", "error", err)
		},
		ConnectUsername: conf.MQTTUsername,
		ConnectPassword: []byte(conf.MQTTPassword),
		WillMessage: &paho.WillMessage{
			QoS:     1,
			Retain:  true,
			Topic:   conf.TopicPrefix + "/status",
			Payload: []byte("offline"),
		},
		ClientConfig: paho.ClientConfig{
			ClientID: conf.TopicPrefix,
			OnClientError: func(err error) {
				slog.Error("Client error", "error", err)
			},
			OnServerDisconnect: func(d *paho.Disconnect) {
				var log *slog.Logger
				if d.Properties != nil {
					log = slog.With("reason", d.Properties.ReasonString)
				} else {
					log = slog.With("reason", d.ReasonCode)
				}
				log.Info("Server requested disconnect")
			},
		},
	}

	client, err := autopaho.NewConnection(context.Background(), cliCfg)
	if err != nil {
		return nil, err
	}

	if err = client.AwaitConnection(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func Disconnect(ctx context.Context, conf *config.Config, client *autopaho.ConnectionManager) error {
	var errs []error

	if _, err := client.Publish(ctx, &paho.Publish{
		Retain:  true,
		Topic:   conf.TopicPrefix + "/status",
		Payload: []byte("offline"),
	}); err != nil {
		errs = append(errs, err)
	}

	if err := client.Disconnect(ctx); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
