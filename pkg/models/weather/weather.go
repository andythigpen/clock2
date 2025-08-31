package weather

import "time"

type WeatherEntity struct {
	State      WeatherCondition        `json:"state"`
	Attributes WeatherEntityAttributes `json:"attributes"`
}

type WeatherEntityAttributes struct {
	Temperature int8  `json:"temperature"`
	Humidity    uint8 `json:"humidity,omitempty"`
}

type ForecastEntity struct {
	Attributes ForecastEntityAttributes `json:"attributes"`
}

type ForecastEntityAttributes struct {
	Forecast []Forecast `json:"forecast"`
}

type Forecast struct {
	DateTime                 time.Time        `json:"datetime"`
	Condition                WeatherCondition `json:"condition"`
	PrecipitationProbability uint8            `json:"precipitation_probability"`
	Temperature              int8             `json:"temperature"`
}

type SunEntity struct {
	Attributes SunEntityAttributes `json:"attributes"`
}

type SunEntityAttributes struct {
	Rising      bool      `json:"rising"`
	NextRising  time.Time `json:"next_rising"`
	NextSetting time.Time `json:"next_setting"`
}
