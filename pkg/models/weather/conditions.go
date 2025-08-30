package weather

import "encoding/json"

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

func (r *WeatherCondition) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
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
	return nil
}
