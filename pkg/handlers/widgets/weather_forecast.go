package widgets

import (
	"context"
	"io"

	"github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/services"
)

type WeatherForecastWidget struct {
	svc *services.HomeAssistantService
}

func NewWeatherForecastWidget(svc *services.HomeAssistantService) Widget {
	return &WeatherForecastWidget{svc}
}

var _ Widget = (*WeatherForecastWidget)(nil)

func (r *WeatherForecastWidget) Render(ctx context.Context, w io.Writer) {
	forecast := r.svc.GetForecast()
	m := view.NewWeatherForecastView(forecast)
	weather.WeatherForecast(m).Render(ctx, w)
}

func (r *WeatherForecastWidget) ShouldDisplay() bool {
	return true
}
