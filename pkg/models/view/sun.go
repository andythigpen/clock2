package view

import (
	"github.com/andythigpen/clock2/pkg/models/weather"
)

type SunView struct {
	Icon         string
	SunDirection string
	Time         string
}

func NewSunView(sun weather.SunEntity) SunView {
	iconName := "sunset"
	direction := "Sunset"
	at := sun.Attributes.NextSetting.Local().Format("03:04")
	if sun.Attributes.Rising {
		iconName = "sunrise"
		direction = "Sunrise"
		at = sun.Attributes.NextRising.Local().Format("03:04")
	}
	return SunView{
		Icon:         AssetIconWeather(iconName),
		SunDirection: direction,
		Time:         at,
	}
}
