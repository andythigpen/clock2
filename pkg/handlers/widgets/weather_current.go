package widgets

import (
	"context"
	"io"

	"github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/services"
)

type WeatherCurrentWidget struct {
	svc *services.HomeAssistantService
}

func NewWeatherCurrentWidget(svc *services.HomeAssistantService) Widget {
	return &WeatherCurrentWidget{svc}
}

var _ Widget = (*WeatherCurrentWidget)(nil)

func (r *WeatherCurrentWidget) Render(ctx context.Context, w io.Writer) {
	current := r.svc.GetWeather()
	forecast := r.svc.GetForecast()
	m := view.NewWeatherCurrentView(current, forecast)
	weather.WeatherCurrent(m).Render(ctx, w)
}

func (r *WeatherCurrentWidget) ShouldDisplay() bool {
	return true
}
