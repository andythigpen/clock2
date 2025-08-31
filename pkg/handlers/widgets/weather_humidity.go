package widgets

import (
	"context"
	"io"

	"github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/services"
)

type WeatherHumidityWidget struct {
	svc *services.HomeAssistantService
}

func NewWeatherHumidityWidget(svc *services.HomeAssistantService) Widget {
	return &WeatherHumidityWidget{svc}
}

var _ Widget = (*WeatherHumidityWidget)(nil)

func (r *WeatherHumidityWidget) Render(ctx context.Context, w io.Writer) {
	current := r.svc.GetWeather()
	m := view.NewWeatherHumidityView(current)
	weather.WeatherHumidity(m).Render(ctx, w)
}

func (r *WeatherHumidityWidget) ShouldDisplay() bool {
	return true
}
