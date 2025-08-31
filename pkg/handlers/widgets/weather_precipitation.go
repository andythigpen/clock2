package widgets

import (
	"context"
	"io"
	"time"

	component "github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/models/weather"
	"github.com/andythigpen/clock2/pkg/services"
)

type WeatherPrecipitationWidget struct {
	svc           *services.HomeAssistantService
	precipitation weather.Forecast
}

func NewWeatherPrecipitationWidget(svc *services.HomeAssistantService) Widget {
	return &WeatherPrecipitationWidget{svc, weather.Forecast{}}
}

var _ Widget = (*WeatherPrecipitationWidget)(nil)

func (r *WeatherPrecipitationWidget) Render(ctx context.Context, w io.Writer) {
	m := view.NewWeatherPrecipitationView(r.precipitation)
	component.WeatherPrecipitation(m).Render(ctx, w)
}

func isPrecipitation(condition weather.WeatherCondition) bool {
	switch condition {
	case weather.Rain, weather.Thunderstorms, weather.ThunderstormsRain, weather.Sleet, weather.Snow, weather.Hail:
		return true
	default:
		return false
	}
}

func (r *WeatherPrecipitationWidget) ShouldDisplay() bool {
	forecast := r.svc.GetForecast()
	count := 0
	for _, hour := range forecast.Attributes.Forecast {
		if hour.DateTime.Before(time.Now()) {
			continue
		}
		count += 1
		if count >= 4 {
			return false
		}
		if !isPrecipitation(hour.Condition) || hour.PrecipitationProbability < 12 {
			continue
		}
		r.precipitation = hour
		return true
	}
	return false
}
