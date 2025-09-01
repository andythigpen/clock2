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

func WeatherConditionIcon(condition weather.WeatherCondition) string {
	var name string
	timeDependent := false
	switch condition {
	case weather.Clear, weather.Cloudy, weather.Fog, weather.PartlyCloudy, weather.Thunderstorms, weather.ThunderstormsRain:
		name = string(condition)
		timeDependent = true
	case weather.Exceptional:
		name = "code-red"
	case weather.Unknown, "":
		name = "code-orange"
	default:
		name = string(condition)
	}
	if timeDependent {
		if isDay() {
			return name + "-day"
		}
		return name + "-night"
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
