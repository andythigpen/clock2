package view

import (
	"strconv"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type WeatherHumidityView struct {
	Icon       string
	Percentage string
}

func NewWeatherHumidityView(current weather.WeatherEntity) WeatherHumidityView {
	return WeatherHumidityView{
		Icon:       AssetIconWeather("humidity", Animated()),
		Percentage: strconv.Itoa(int(current.Attributes.Humidity)),
	}
}
