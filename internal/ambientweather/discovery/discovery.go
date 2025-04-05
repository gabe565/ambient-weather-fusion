package discovery

import (
	"gabe565.com/ambient-weather-fusion/internal/config"
	"k8s.io/utils/ptr"
)

func NewPayload(conf *config.Config, version string) Payload { //nolint:funlen
	components := map[Topic]Component{
		TopicTemperature: {
			Platform:                  PlatformSensor,
			UnitOfMeasurement:         UnitFahrenheit,
			DeviceClass:               DeviceClassTemperature,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
		TopicHumidity: {
			Platform:                  PlatformSensor,
			UnitOfMeasurement:         UnitPercent,
			DeviceClass:               DeviceClassHumidity,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
		TopicWindSpeed: {
			Platform:                  PlatformSensor,
			UnitOfMeasurement:         UnitMPH,
			DeviceClass:               DeviceClassWindSpeed,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
		TopicWindGust: {
			Platform:                  PlatformSensor,
			Name:                      "Wind gust",
			UnitOfMeasurement:         UnitMPH,
			DeviceClass:               DeviceClassWindSpeed,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
		TopicMaxDailyGust: {
			Platform:                  PlatformSensor,
			Name:                      "Max daily gust",
			UnitOfMeasurement:         UnitMPH,
			DeviceClass:               DeviceClassWindSpeed,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
		TopicUVIndex: {
			Platform:                  PlatformSensor,
			Name:                      "UV index",
			UnitOfMeasurement:         UnitIndex,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
		TopicSolarRadiation: {
			Platform:                  PlatformSensor,
			UnitOfMeasurement:         UnitWattsPerSqMeter,
			DeviceClass:               DeviceClassIrradiance,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
			EnabledByDefault:          ptr.To(false),
		},
		TopicHourlyRain: {
			Platform:                  PlatformSensor,
			Name:                      "Hourly rain",
			UnitOfMeasurement:         UnitInchesPerHour,
			DeviceClass:               DeviceClassPrecipitationIntensity,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 2,
		},
		TopicDailyRain: {
			Platform:                  PlatformSensor,
			Name:                      "Daily rain",
			UnitOfMeasurement:         UnitInches,
			DeviceClass:               DeviceClassPrecipitation,
			StateClass:                StateClassTotal,
			SuggestedDisplayPrecision: 2,
		},
		TopicWeeklyRain: {
			Platform:                  PlatformSensor,
			Name:                      "Weekly rain",
			UnitOfMeasurement:         UnitInches,
			DeviceClass:               DeviceClassPrecipitation,
			StateClass:                StateClassTotal,
			SuggestedDisplayPrecision: 2,
			EnabledByDefault:          ptr.To(false),
		},
		TopicMonthlyRain: {
			Platform:                  PlatformSensor,
			Name:                      "Monthly rain",
			UnitOfMeasurement:         UnitInches,
			DeviceClass:               DeviceClassPrecipitation,
			StateClass:                StateClassTotal,
			SuggestedDisplayPrecision: 2,
			EnabledByDefault:          ptr.To(false),
		},
		TopicRelativePressure: {
			Platform:                  PlatformSensor,
			Name:                      "Relative pressure",
			UnitOfMeasurement:         UnitInHg,
			DeviceClass:               DeviceClassPressure,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 2,
		},
		TopicAbsolutePressure: {
			Platform:                  PlatformSensor,
			Name:                      "Absolute pressure",
			UnitOfMeasurement:         UnitInHg,
			DeviceClass:               DeviceClassPressure,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 2,
			EnabledByDefault:          ptr.To(false),
		},
		TopicLastRain: {
			Platform:         PlatformSensor,
			Name:             "Last rain",
			DeviceClass:      DeviceClassTimestamp,
			EnabledByDefault: ptr.To(false),
			Icon:             "mdi:water",
		},
		TopicFeelsLike: {
			Platform:                  PlatformSensor,
			Name:                      "Feels like",
			UnitOfMeasurement:         UnitFahrenheit,
			DeviceClass:               DeviceClassTemperature,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
		TopicDewPoint: {
			Platform:                  PlatformSensor,
			Name:                      "Dew point",
			UnitOfMeasurement:         UnitFahrenheit,
			DeviceClass:               DeviceClassTemperature,
			StateClass:                StateClassMeasurement,
			SuggestedDisplayPrecision: 1,
		},
	}

	for topic, sensor := range components {
		sensor.ObjectID = conf.BaseTopic + "_" + string(topic)
		sensor.UniqueID = conf.BaseTopic + "_" + string(topic)
		sensor.ValueTemplate = "{{ value_json." + string(topic) + " }}"
		components[topic] = sensor
	}

	return Payload{
		AvailabilityTopic: conf.BaseTopic + "/status",
		Device: Device{
			Identifiers: conf.BaseTopic,
			Name:        conf.HADeviceName,
			SWVersion:   version,
		},
		Origin: Origin{
			Name:       "Ambient Weather Fusion",
			SWVersion:  version,
			SupportURL: "https://github.com/gabe565/ambient-weather-fusion",
		},
		StateTopic: conf.BaseTopic,
		Components: components,
	}
}
