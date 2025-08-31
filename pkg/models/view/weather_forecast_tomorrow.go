package view

import (
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type WeatherForecastTomorrowView struct {
	Icon          string
	TemperatureLo string
	TemperatureHi string
}

func NewWeatherForecastTomorrowView(forecast weather.ForecastEntity) WeatherForecastTomorrowView {
	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	dayAfterTomorrow := time.Now().AddDate(0, 0, 2).Truncate(24 * time.Hour)
	conditions := map[weather.WeatherCondition]int{}
	tempHi := int8(-127)
	tempLo := int8(127)
	for _, hour := range forecast.Attributes.Forecast {
		if !hour.DateTime.After(tomorrow) {
			continue
		}
		if !hour.DateTime.Before(dayAfterTomorrow) {
			break
		}
		conditions[hour.Condition] += 1
		if hour.Temperature < tempLo {
			tempLo = hour.Temperature
		} else if hour.Temperature > tempHi {
			tempHi = hour.Temperature
		}
	}
	condition := weather.Unknown
	maximum := 0
	for c, num := range conditions {
		if num > maximum {
			condition = c
			maximum = num
		}
	}
	return WeatherForecastTomorrowView{
		Icon:          AssetIconWeather(WeatherConditionIcon(condition), Animated()),
		TemperatureLo: strconv.Itoa(int(tempLo)),
		TemperatureHi: strconv.Itoa(int(tempHi)),
	}
}
