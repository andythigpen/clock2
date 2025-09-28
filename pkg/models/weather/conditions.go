package weather

import (
	"encoding/json"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type WeatherCondition string

const (
	Clear             WeatherCondition = "clear"
	Cloudy            WeatherCondition = "cloudy"
	Exceptional       WeatherCondition = "exceptional"
	Fog               WeatherCondition = "fog"
	Hail              WeatherCondition = "hail"
	PartlyCloudy      WeatherCondition = "partly-cloudy"
	Rain              WeatherCondition = "rain"
	Sleet             WeatherCondition = "sleet"
	Snow              WeatherCondition = "snow"
	Thunderstorms     WeatherCondition = "thunderstorms"
	ThunderstormsRain WeatherCondition = "thunderstorms-rain"
	Unknown           WeatherCondition = "unknown"
	Windy             WeatherCondition = "windy"
)

var AllConditions = []WeatherCondition{
	Clear,
	Cloudy,
	Exceptional,
	Fog,
	Hail,
	PartlyCloudy,
	Rain,
	Sleet,
	Snow,
	Thunderstorms,
	ThunderstormsRain,
	Unknown,
	Windy,
}

func (r WeatherCondition) String() string {
	switch r {
	case Exceptional:
		return "Warning"
	case PartlyCloudy:
		return "Partly Cloudy"
	case ThunderstormsRain:
		return "Thunderstorms"
	default:
		caser := cases.Title(language.AmericanEnglish)
		return caser.String(string(r))
	}
}

func (r *WeatherCondition) FromString(s string) {
	switch s {
	case "clear-day", "clear-night", "sunny":
		*r = Clear
	case "cloudy":
		*r = Cloudy
	case "fog":
		*r = Fog
	case "hail":
		*r = Hail
	case "lightning", "thunderstorms":
		*r = Thunderstorms
	case "lightning-rainy", "thunderstorms-rain":
		*r = ThunderstormsRain
	case "partly-cloudy", "partlycloudy":
		*r = PartlyCloudy
	case "pouring", "rainy", "rain":
		*r = Rain
	case "snowy", "snow":
		*r = Snow
	case "snowy-rainy", "snowy-rain", "sleet":
		*r = Sleet
	case "windy", "windy-variant":
		*r = Windy
	case "exceptional", "alert":
		*r = Exceptional
	default:
		*r = Unknown
	}
}

func (r *WeatherCondition) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	r.FromString(s)
	return nil
}
