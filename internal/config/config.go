package config

import (
	"net/url"
	"strconv"
	"time"

	"gabe565.com/ambient-weather-fusion/internal/location"
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
		Radius:          4,
		Limit:           100,
		MaxReadingAge:   10 * time.Minute,
		DiscoveryPrefix: "homeassistant",
		TopicPrefix:     "ambient_weather_fusion",
		DeviceName:      "Ambient Weather Fusion",
	}
}

func (c *Config) BuildURL() *url.URL {
	u := *c.RequestURL.URL
	q := u.Query()
	lat1, lon1 := location.Shift(c.Latitude, c.Longitude, -c.Radius, -c.Radius)
	lat2, lon2 := location.Shift(c.Latitude, c.Longitude, c.Radius, c.Radius)
	q.Set("$publicBox[0][0]", strconv.FormatFloat(lon1, 'f', -1, 64))
	q.Set("$publicBox[0][1]", strconv.FormatFloat(lat1, 'f', -1, 64))
	q.Set("$publicBox[1][0]", strconv.FormatFloat(lon2, 'f', -1, 64))
	q.Set("$publicBox[1][1]", strconv.FormatFloat(lat2, 'f', -1, 64))
	q.Set("$limit", strconv.Itoa(c.Limit))
	u.RawQuery = q.Encode()
	return &u
}
