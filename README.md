# Ambient Weather Fusion

[![Build](https://github.com/gabe565/ambient-weather-fusion/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/ambient-weather-fusion/actions/workflows/build.yml)

A tool for aggregating Ambient Weather Network data to provide more reliable local readings. Instead of using a single station, this tool calculates the median values from nearby sensors to reduce bad readings.

## Features

- Designed for integration with Home Assistant through MQTT
- Reports values for temperature, humidity, wind speed, pressure, and more
- Smooth out bad readings from individual weather stations

## Usage

Builds are published as a Docker container to [ghcr.io/gabe565/ambient-weather-fusion](https://ghcr.io/gabe565/ambient-weather-fusion).

```shell
docker run -d ghcr.io/gabe565/ambient-weather-fusion \
  --mqtt-url mqtt://localhost:1883 \
  --latitude 35.4689 \
  --longitude -97.5195
```

- [Command line reference](docs/ambient-weather-fusion.md)
- [Environment variable reference](docs/envs.md)

## Sensors

<p align="center">
  <img alt="Ambient Weather Fusion sensors screenshot" width="250" src="https://github.com/user-attachments/assets/16540f7b-a896-40f0-88a4-18b9ee107126#gh-dark-mode-only">
  <img alt="Ambient Weather Fusion sensors screenshot" width="250" src="https://github.com/user-attachments/assets/7a450815-61b9-4c49-a6f8-f83bc311760f#gh-light-mode-only">
</p>
