package view

import (
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

var gradient = []string{
	// NOTE: the from and to comments are necessary so the styles are included in style.css
	"#241b63", // <30 from-[#241b63] to-[#241b63]
	"#20348d", // 30 from-[#20348d] to-[#20348d]
	"#224396", // 35 from-[#224396] to-[#224396]
	"#3d64ae", // 40 from-[#3d64ae] to-[#3d64ae]
	"#5583c1", // 45 from-[#5583c1] to-[#5583c1]
	"#89addc", // 50 from-[#89addc] to-[#89addc]
	"#9bbae2", // 55 from-[#9bbae2] to-[#9bbae2]
	"#9acdcf", // 60 from-[#9acdcf] to-[#9acdcf]
	"#a2cfa4", // 65 from-[#a2cfa4] to-[#a2cfa4]
	"#d7de7e", // 70 from-[#d7de7e] to-[#d7de7e]
	"#f4d862", // 75 from-[#f4d862] to-[#f4d862]
	"#f29400", // 80 from-[#f29400] to-[#f29400]
	"#eb561e", // 85 from-[#eb561e] to-[#eb561e]
	"#c30507", // 90 from-[#c30507] to-[#c30507]
	"#700318", // >90 from-[#700318] to-[#700318]
}

type WeatherForecastTomorrowView struct {
	Icon          string
	TemperatureLo string
	TemperatureHi string
	ColorLo       string
	ColorHi       string
}

func getGradientColor(temp int8) string {
	if temp < 30 {
		return gradient[0]
	}
	if temp > 90 {
		return gradient[len(gradient)-1]
	}
	// output = output_start + ((output_end - output_start) / (input_end - input_start)) * (input - input_start)
	idx := min(1+int((float64(len(gradient)-3)/float64(90-30))*float64(int(temp)-30)), len(gradient)-1)
	return gradient[idx]
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
		ColorLo:       getGradientColor(tempLo),
		ColorHi:       getGradientColor(tempHi),
	}
}
