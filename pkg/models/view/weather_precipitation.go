package view

import (
	"strconv"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type WeatherPrecipitationView struct {
	Icon              string
	Hour              string
	AmPm              string
	Percent           string
	PrecipitationType string
}

func NewWeatherPrecipitationView(forecast weather.Forecast) WeatherPrecipitationView {
	ampm := "AM"
	if forecast.DateTime.Hour() >= 12 {
		ampm = "PM"
	}
	return WeatherPrecipitationView{
		Icon:              AssetIconWeather(WeatherConditionIcon(forecast.Condition), Animated()),
		Hour:              forecast.DateTime.Format("03"),
		AmPm:              ampm,
		Percent:           strconv.Itoa(int(forecast.PrecipitationProbability)),
		PrecipitationType: forecast.Condition.String(),
	}
}
