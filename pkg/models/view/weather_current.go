package view

import (
	"fmt"
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type WeatherCurrentView struct {
	IconSrc                  string
	Temperature              string
	TemperatureDirectionIcon string
	DisplayDirection         bool
}

func NewWeatherCurrentView(weather weather.WeatherEntity, forecast weather.ForecastEntity) WeatherCurrentView {
	currentTemperature := weather.Attributes.Temperature
	now := time.Now()

	direction := ""
	displayDirection := false
	for _, f := range forecast.Attributes.Forecast {
		if f.DateTime.Before(now) {
			continue
		}
		if f.Temperature > currentTemperature {
			direction = "high"
			displayDirection = f.Temperature-currentTemperature > 2
		} else if f.Temperature < currentTemperature {
			direction = "low"
			displayDirection = currentTemperature-f.Temperature > 2
		}
		break
	}
	return WeatherCurrentView{
		IconSrc:                  AssetIconWeather(WeatherConditionIcon(weather.State)),
		Temperature:              strconv.Itoa(int(currentTemperature)),
		TemperatureDirectionIcon: AssetIconWeather(fmt.Sprintf("pressure-%s", direction), Animated()),
		DisplayDirection:         displayDirection,
	}
}
