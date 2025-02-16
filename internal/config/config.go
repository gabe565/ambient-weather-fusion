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

	DiscoveryPrefix string
	TopicPrefix     string
	DeviceName      string

	MQTTURL                pflagx.URL
	MQTTUsername           string
	MQTTPassword           string
	MQTTInsecureSkipVerify bool
	MQTTKeepAlive          uint16
	MQTTSessionExpiry      uint32
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
		Radius:            4,
		Limit:             100,
		MaxReadingAge:     10 * time.Minute,
		DiscoveryPrefix:   "homeassistant",
		TopicPrefix:       "ambient_weather_fusion",
		DeviceName:        "Ambient Weather Fusion",
		MQTTKeepAlive:     60,
		MQTTSessionExpiry: 60,
	}
}
