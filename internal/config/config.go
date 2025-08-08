package config

import (
	"net/url"
	"time"

	"gabe565.com/utils/pflagx"
)

type Config struct {
	RequestURL    pflagx.URL
	Latitude      float64
	Longitude     float64
	Radius        float64
	Limit         int
	MaxReadingAge time.Duration

	MQTTURL                pflagx.URL
	MQTTUsername           string
	MQTTPassword           string
	MQTTCAPath             string
	MQTTClientCertPath     string
	MQTTClientKeyPath      string
	MQTTInsecureSkipVerify bool
	MQTTKeepAlive          uint16
	MQTTSessionExpiry      uint32

	BaseTopic        string
	HADiscoveryTopic string
	HAStatusTopic    string
	HADeviceName     string
}

func New() *Config {
	return &Config{
		RequestURL: pflagx.URL{
			URL: &url.URL{
				Scheme: "https",
				Host:   "lightning.ambientweather.net",
				Path:   "/devices",
			},
		},
		Radius:        4,
		Limit:         100,
		MaxReadingAge: 10 * time.Minute,

		MQTTKeepAlive:     60,
		MQTTSessionExpiry: 60,

		BaseTopic:        "ambient_weather_fusion",
		HADiscoveryTopic: "homeassistant",
		HAStatusTopic:    "homeassistant/status",
		HADeviceName:     "Ambient Weather Fusion",
	}
}
