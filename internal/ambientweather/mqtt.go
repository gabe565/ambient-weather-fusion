package ambientweather

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"net/url"
	"os"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func (s *Server) ConnectMQTT(ctx context.Context) error {
	var pool *x509.CertPool
	if s.conf.MQTTCAPath != "" {
		var err error
		if pool, err = loadCACert(s.conf.MQTTCAPath); err != nil {
			return err
		}
	}

	log := slog.With("url", s.conf.MQTTURL.String())
	cliCfg := autopaho.ClientConfig{
		ServerUrls: []*url.URL{s.conf.MQTTURL.URL},
		TlsCfg: &tls.Config{
			InsecureSkipVerify: s.conf.MQTTInsecureSkipVerify, //nolint:gosec
			RootCAs:            pool,
		},
		KeepAlive:             s.conf.MQTTKeepAlive,
		SessionExpiryInterval: s.conf.MQTTSessionExpiry,
		OnConnectionUp: func(_ *autopaho.ConnectionManager, _ *paho.Connack) {
			log.Info("Connected to MQTT")
			if err := s.PublishStatus(ctx, true); err != nil {
				log.Error("Failed to publish status message", "error", err)
			}
			if s.conf.HAStatusTopic != "" {
				if _, err := s.mqtt.Subscribe(ctx, &paho.Subscribe{
					Subscriptions: []paho.SubscribeOptions{
						{Topic: s.conf.HAStatusTopic, QoS: 1},
					},
				}); err != nil {
					slog.Error("Failed to subscribe to Home Assistant status topic", "error", err)
				}
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
			Topic:   s.conf.BaseTopic + "/status",
			Payload: []byte("offline"),
		},
		ClientConfig: paho.ClientConfig{
			ClientID: s.conf.BaseTopic,
			OnPublishReceived: []func(received paho.PublishReceived) (bool, error){
				func(r paho.PublishReceived) (bool, error) {
					if r.Packet.Topic == s.conf.HAStatusTopic && string(r.Packet.Payload) == "online" {
						if s.lastPayload == nil {
							return true, nil
						}
						return true, s.PublishData(ctx, s.lastPayload)
					}
					return false, nil
				},
			},
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

func loadCACert(path string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	pemCerts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for len(pemCerts) != 0 {
		var block *pem.Block
		if block, pemCerts = pem.Decode(pemCerts); block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		pool.AddCert(cert)
	}

	return pool, nil
}
