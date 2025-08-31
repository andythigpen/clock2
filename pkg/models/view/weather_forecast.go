package view

import (
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type ForecastHour struct {
	Hour        string
	AmPm        string
	Icon        string
	Temperature string
}

type WeatherForecastView struct {
	Hours []ForecastHour
}

func NewWeatherForecastView(forecast weather.ForecastEntity) WeatherForecastView {
	hours := []ForecastHour{}
	for _, hour := range forecast.Attributes.Forecast {
		if hour.DateTime.After(time.Now()) {
			ampm := "AM"
			if hour.DateTime.Hour() >= 12 {
				ampm = "PM"
			}
			hours = append(hours, ForecastHour{
				Hour:        hour.DateTime.Format("03"),
				AmPm:        ampm,
				Icon:        AssetIconWeather(WeatherConditionIcon(hour.Condition)),
				Temperature: strconv.Itoa(int(hour.Temperature)),
			})
			if len(hours) >= 3 {
				break
			}
		}
	}
	return WeatherForecastView{Hours: hours}
}
