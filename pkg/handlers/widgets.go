package handlers

import (
	"context"
	"io"

	. "github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/services"
)

type Widget interface {
	ShouldDisplay() bool
	Render(context.Context, io.Writer)
}

type WeatherCurrentWidget struct {
	svc *services.HomeAssistantService
}

var _ Widget = (*WeatherCurrentWidget)(nil)

func (r *WeatherCurrentWidget) Render(ctx context.Context, w io.Writer) {
	weather := r.svc.GetWeather()
	forecast := r.svc.GetForecast()
	m := view.NewWeatherCurrentView(weather, forecast)
	WeatherCurrent(m).Render(ctx, w)
}

func (r *WeatherCurrentWidget) ShouldDisplay() bool {
	return true
}

type WeatherForecastWidget struct {
	svc *services.HomeAssistantService
}

var _ Widget = (*WeatherForecastWidget)(nil)

func (r *WeatherForecastWidget) Render(ctx context.Context, w io.Writer) {
	forecast := r.svc.GetForecast()
	m := view.NewWeatherForecastView(forecast)
	WeatherForecast(m).Render(ctx, w)
}

func (r *WeatherForecastWidget) ShouldDisplay() bool {
	return true
}

type WeatherForecastTomorrowWidget struct {
	svc *services.HomeAssistantService
}

var _ Widget = (*WeatherForecastTomorrowWidget)(nil)

func (r *WeatherForecastTomorrowWidget) Render(ctx context.Context, w io.Writer) {
	forecast := r.svc.GetForecast()
	m := view.NewWeatherForecastTomorrowView(forecast)
	WeatherForecastTomorrow(m).Render(context.Background(), w)
}

func (r *WeatherForecastTomorrowWidget) ShouldDisplay() bool {
	return true
}
