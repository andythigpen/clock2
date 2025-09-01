package view

import (
	"strconv"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type WeatherPrecipitationView struct {
	Icon    string
	Hour    string
	Percent string
}

func NewWeatherPrecipitationView(forecast weather.Forecast) WeatherPrecipitationView {
	return WeatherPrecipitationView{
		Icon:    AssetIconWeather(WeatherConditionIcon(forecast.Condition), Animated()),
		Hour:    forecast.DateTime.Format("03"),
		Percent: strconv.Itoa(int(forecast.PrecipitationProbability)),
	}
}
