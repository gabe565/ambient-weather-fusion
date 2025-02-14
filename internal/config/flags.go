package config

import "github.com/spf13/cobra"

const (
	FlagRequestURL    = "request-url"
	FlagLatitude      = "latitude"
	FlagLongitude     = "longitude"
	FlagRadius        = "radius"
	FlagMaxReadingAge = "max-reading-age"

	FlagDiscoveryPrefix = "discovery-prefix"
	FlagTopicPrefix     = "topic-prefix"
	FlagDeviceName      = "device-name"

	FlagMQTTURL      = "mqtt-url"
	FlagMQTTUsername = "mqtt-username"
	FlagMQTTPassword = "mqtt-password"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.Var(&c.RequestURL, FlagRequestURL, "Ambient Weather API URL")
	fs.Float64Var(&c.Latitude, FlagLatitude, c.Latitude, "Latitude of center")
	fs.Float64Var(&c.Longitude, FlagLongitude, c.Longitude, "Longitude of center")
	fs.Float64Var(&c.Radius, FlagRadius, c.Radius, "Radius in miles")
	fs.DurationVar(&c.MaxReadingAge, FlagMaxReadingAge, c.MaxReadingAge, "Maximum age of a reading to be included")

	fs.StringVar(&c.DiscoveryPrefix, FlagDiscoveryPrefix, c.DiscoveryPrefix, "Home Assistant discovery prefix")
	fs.StringVar(&c.TopicPrefix, FlagTopicPrefix, c.TopicPrefix, "Topic prefix")
	fs.StringVar(&c.DeviceName, FlagDeviceName, c.DeviceName, "Name of the device to add to Home Assistant")

	fs.Var(&c.MQTTURL, FlagMQTTURL, "MQTT server URL")
	fs.StringVar(&c.MQTTUsername, FlagMQTTUsername, c.MQTTUsername, "MQTT username")
	fs.StringVar(&c.MQTTPassword, FlagMQTTPassword, c.MQTTPassword, "MQTT password")
}
