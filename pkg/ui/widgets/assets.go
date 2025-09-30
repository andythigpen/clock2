package widgets

import (
	"flag"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
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

type weatherIconOptions struct {
	iconType   string
	size       int
	isAnimated bool
}

type iconOption func(*weatherIconOptions)

func Animated() iconOption {
	return func(o *weatherIconOptions) {
		o.isAnimated = true
	}
}

func WithSize(size int) iconOption {
	return func(o *weatherIconOptions) {
		o.size = size
	}
}

func getAssetIconPath(name string, options ...iconOption) string {
	o := &weatherIconOptions{
		iconType:   "weather",
		size:       480,
		isAnimated: false,
	}
	for _, opt := range options {
		opt(o)
	}
	if o.isAnimated {
		return path.Join("assets/icons", o.iconType, "animated", strconv.Itoa(o.size), fmt.Sprintf("%s.gif", name))
	}
	return path.Join("assets/icons", o.iconType, "static", strconv.Itoa(o.size), fmt.Sprintf("%s.png", name))
}

func getWeatherConditionIconName(condition weather.WeatherCondition) string {
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
