package widgets

import (
	"context"
	"io"
	"time"

	"github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/services"
)

type WeatherForecastTomorrowWidget struct {
	svc *services.HomeAssistantService
}

func NewWeatherForecastTomorrowWidget(svc *services.HomeAssistantService) Widget {
	return &WeatherForecastTomorrowWidget{svc}
}

var _ Widget = (*WeatherForecastTomorrowWidget)(nil)

func (r *WeatherForecastTomorrowWidget) Render(ctx context.Context, w io.Writer) {
	forecast := r.svc.GetForecast()
	m := view.NewWeatherForecastTomorrowView(forecast)
	weather.WeatherForecastTomorrow(m).Render(ctx, w)
}

func (r *WeatherForecastTomorrowWidget) ShouldDisplay() bool {
	return time.Now().Hour() >= 17
}
