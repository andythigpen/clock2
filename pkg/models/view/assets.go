package view

import (
	"fmt"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type weatherIconOpt struct {
	animated bool
}
type WeatherIconOption func(*weatherIconOpt)

func Animated() WeatherIconOption {
	return func(w *weatherIconOpt) {
		w.animated = true
	}
}

func isDay() bool {
	now := time.Now()
	hour := now.Hour()
	return hour >= 6 && hour <= 19
}

func WeatherIconName(condition weather.WeatherCondition) string {
	name := "unknown"
	// TODO: just rename the files to match the conditions
	switch condition {
	case weather.Clear:
		if isDay() {
			name = "clear-day"
		} else {
			name = "clear-night"
		}
	case weather.Cloudy:
		if isDay() {
			name = "overcast-day"
		} else {
			name = "overcast-night"
		}
	case weather.Exceptional:
		name = "code-red"
	case weather.Fog:
		name = "fog"
	case weather.PartlyCloudy:
		if isDay() {
			name = "partly-cloudy-day"
		} else {
			name = "partly-cloudy-night"
		}
	case weather.Rain:
		name = "rain"
	case weather.Sleet:
		name = "sleet"
	case weather.Snow:
		name = "snow"
	case weather.Thunderstorms:
		if isDay() {
			name = "thunderstorms-day"
		} else {
			name = "thunderstorms-night"
		}
	case weather.ThunderstormsRain:
		if isDay() {
			name = "thunderstorms-day-rain"
		} else {
			name = "thunderstorms-night-rain"
		}
	case weather.Windy:
		name = "wind"
	}
	return name
}

func AssetIconWeather(name string, opts ...WeatherIconOption) string {
	options := weatherIconOpt{}
	for _, o := range opts {
		o(&options)
	}
	animated := "static"
	if options.animated {
		animated = "animated"
	}
	return fmt.Sprintf("/assets/icons/weather/%s/%s.svg", animated, name)
}
