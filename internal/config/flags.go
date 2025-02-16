package config

import "github.com/spf13/cobra"

const (
	FlagRequestURL    = "request-url"
	FlagLatitude      = "latitude"
	FlagLongitude     = "longitude"
	FlagRadius        = "radius"
	FlagMaxReadingAge = "max-reading-age"

	FlagMQTTURL           = "mqtt-url"
	FlagMQTTUsername      = "mqtt-username"
	FlagMQTTPassword      = "mqtt-password"
	FlagMQTTInsecure      = "mqtt-insecure"
	FlagMQTTKeepAlive     = "mqtt-keep-alive"
	FlagMQTTSessionExpiry = "mqtt-session-expiry"

	FlagBaseTopic        = "base-topic"
	FlagHADiscoveryTopic = "ha-discovery-topic"
	FlagHAStatusTopic    = "ha-status-topic"
	FlagHADeviceName     = "ha-device-name"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.Var(&c.RequestURL, FlagRequestURL, "Ambient Weather API URL")
	fs.Float64Var(&c.Latitude, FlagLatitude, c.Latitude, "Latitude of center")
	fs.Float64Var(&c.Longitude, FlagLongitude, c.Longitude, "Longitude of center")
	fs.Float64Var(&c.Radius, FlagRadius, c.Radius, "Radius in miles")
	fs.DurationVar(&c.MaxReadingAge, FlagMaxReadingAge, c.MaxReadingAge, "Maximum age of a reading to be included")

	fs.Var(&c.MQTTURL, FlagMQTTURL, "MQTT server URL")
	fs.StringVar(&c.MQTTUsername, FlagMQTTUsername, c.MQTTUsername, "MQTT username")
	fs.StringVar(&c.MQTTPassword, FlagMQTTPassword, c.MQTTPassword, "MQTT password")
	fs.BoolVar(&c.MQTTInsecureSkipVerify, FlagMQTTInsecure, c.MQTTInsecureSkipVerify, "Skip MQTT TLS verification")
	fs.Uint16Var(&c.MQTTKeepAlive, FlagMQTTKeepAlive, c.MQTTKeepAlive, "MQTT keep alive interval in seconds")
	fs.Uint32Var(&c.MQTTSessionExpiry, FlagMQTTSessionExpiry, c.MQTTSessionExpiry, "MQTT session expiry interval in seconds")

	fs.StringVar(&c.BaseTopic, FlagBaseTopic, c.BaseTopic, "MQTT base topic")
	fs.StringVar(&c.HADiscoveryTopic, FlagHADiscoveryTopic, c.HADiscoveryTopic, "Home Assistant discovery topic")
	fs.StringVar(&c.HAStatusTopic, FlagHAStatusTopic, c.HAStatusTopic, "Home Assistant status topic")
	fs.StringVar(&c.HADeviceName, FlagHADeviceName, c.HADeviceName, "Name of the device to add to Home Assistant")
}
