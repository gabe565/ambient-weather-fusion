## ambient-weather-fusion

Integrate consensus-based Ambient Weather readings into Home Assistant

```
ambient-weather-fusion [flags]
```

### Options

```
      --device-name string         Name of the device to add to Home Assistant (default "Ambient Weather Fusion")
      --discovery-prefix string    Home Assistant discovery prefix (default "homeassistant")
  -h, --help                       help for ambient-weather-fusion
      --latitude float             Latitude of center
      --longitude float            Longitude of center
      --max-reading-age duration   Maximum age of a reading to be included (default 10m0s)
      --mqtt-password string       MQTT password
      --mqtt-url string            MQTT server URL
      --mqtt-username string       MQTT username
      --radius float               Radius in miles (default 4)
      --request-url string         Ambient Weather API URL (default "https://lightning.ambientweather.net/devices")
      --topic-prefix string        Topic prefix (default "ambient_weather_fusion")
  -v, --version                    version for ambient-weather-fusion
```

