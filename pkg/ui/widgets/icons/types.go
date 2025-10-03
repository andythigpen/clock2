package icons

import (
	"flag"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type IconType string

const (
	IconClearDay               IconType = "clear-day"
	IconClearNight             IconType = "clear-night"
	IconCloudyDay              IconType = "cloudy-day"
	IconCloudyNight            IconType = "cloudy-night"
	IconCodeOrange             IconType = "code-orange"
	IconCodeRed                IconType = "code-red"
	IconFogDay                 IconType = "fog-day"
	IconFogNight               IconType = "fog-night"
	IconHail                   IconType = "hail"
	IconHumidity               IconType = "humidity"
	IconPartlyCloudyDay        IconType = "partly-cloudy-day"
	IconPartlyCloudyNight      IconType = "partly-cloudy-night"
	IconPressureHigh           IconType = "pressure-high"
	IconPressureLow            IconType = "pressure-low"
	IconRain                   IconType = "rain"
	IconSleet                  IconType = "sleet"
	IconSnow                   IconType = "snow"
	IconSunrise                IconType = "sunrise"
	IconSunset                 IconType = "sunset"
	IconThunderstormsDay       IconType = "thunderstorms-day"
	IconThunderstormsNight     IconType = "thunderstorms-night"
	IconThunderstormsRainDay   IconType = "thunderstorms-rain-day"
	IconThunderstormsRainNight IconType = "thunderstorms-rain-night"
	IconWindy                  IconType = "windy"
)

var (
	uiTestDayNight = flag.String("ui-test-day-night", "", "set to day or night to change asset icons")
)

func isDay() bool {
	if len(*uiTestDayNight) > 0 {
		return *uiTestDayNight == "day"
	}
	now := time.Now()
	hour := now.Hour()
	return hour >= 6 && hour <= 19
}

func GetWeatherConditionIconType(condition weather.WeatherCondition) IconType {
	var name IconType
	switch condition {
	case weather.Clear, weather.Cloudy, weather.Fog, weather.PartlyCloudy, weather.Thunderstorms, weather.ThunderstormsRain:
		name = IconType(condition)
		if isDay() {
			return name + "-day"
		}
		return name + "-night"
	case weather.Exceptional:
		name = "code-red"
	case weather.Unknown, "":
		name = "code-orange"
	default:
		name = IconType(condition)
	}
	return name
}

type weatherIconOptions struct {
	iconType string
	size     int
}

type iconOption func(*weatherIconOptions)

func WithSize(size int) iconOption {
	return func(o *weatherIconOptions) {
		o.size = size
	}
}

func GetStaticIconPath(iconType IconType, options ...iconOption) string {
	o := &weatherIconOptions{
		iconType: "weather",
		size:     480,
	}
	for _, opt := range options {
		opt(o)
	}
	return path.Join("assets/icons", o.iconType, "static", strconv.Itoa(o.size), fmt.Sprintf("%s.png", string(iconType)))
}
