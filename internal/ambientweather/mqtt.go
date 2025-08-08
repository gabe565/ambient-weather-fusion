package ambientweather

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"gabe565.com/ambient-weather-fusion/internal/config"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func (s *Server) ConnectMQTT(ctx context.Context) error {
	tlsConf, err := newMQTTTLSConfig(s.conf)
	if err != nil {
		return fmt.Errorf("failed to create TLS config: %w", err)
	}

	log := slog.With("url", s.conf.MQTTURL.String())
	cliCfg := autopaho.ClientConfig{
		ServerUrls:            []*url.URL{s.conf.MQTTURL.URL},
		TlsCfg:                tlsConf,
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

	if s.mqtt, err = autopaho.NewConnection(context.Background(), cliCfg); err != nil {
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

func newMQTTTLSConfig(conf *config.Config) (*tls.Config, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: conf.MQTTInsecureSkipVerify, //nolint:gosec
	}

	if conf.MQTTCAPath != "" {
		pool, err := loadCACert(conf.MQTTCAPath)
		if err != nil {
			return nil, err
		}

		tlsConf.RootCAs = pool
	}

	if conf.MQTTClientCertPath != "" || conf.MQTTClientKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(conf.MQTTClientCertPath, conf.MQTTClientKeyPath)
		if err != nil {
			return nil, err
		}

		tlsConf.Certificates = []tls.Certificate{cert}
	}

	return tlsConf, nil
}
