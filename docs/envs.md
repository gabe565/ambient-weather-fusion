# Environment Variables

| Name | Usage | Default |
| --- | --- | --- |
| `AW_BASE_TOPIC` | MQTT base topic | `ambient_weather_fusion` |
| `AW_HA_DEVICE_NAME` | Name of the device to add to Home Assistant | `Ambient Weather Fusion` |
| `AW_HA_DISCOVERY_TOPIC` | Home Assistant discovery topic | `homeassistant` |
| `AW_HA_STATUS_TOPIC` | Home Assistant status topic | `homeassistant/status` |
| `AW_LATITUDE` | Latitude of center | `0` |
| `AW_LONGITUDE` | Longitude of center | `0` |
| `AW_MAX_READING_AGE` | Maximum age of a reading to be included | `10m0s` |
| `AW_MQTT_CA` | MQTT CA certificate file path | ` ` |
| `AW_MQTT_INSECURE` | Skip MQTT TLS verification | `false` |
| `AW_MQTT_KEEP_ALIVE` | MQTT keep alive interval in seconds | `60` |
| `AW_MQTT_PASSWORD` | MQTT password | ` ` |
| `AW_MQTT_SESSION_EXPIRY` | MQTT session expiry interval in seconds | `60` |
| `AW_MQTT_URL` | MQTT server URL | ` ` |
| `AW_MQTT_USERNAME` | MQTT username | ` ` |
| `AW_RADIUS` | Radius in miles | `4` |
| `AW_REQUEST_URL` | Ambient Weather API URL | `https://lightning.ambientweather.net/devices` |