package handlers

import (
	"context"
	"io"
	"time"

	. "github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/models/weather"
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
	WeatherForecastTomorrow(m).Render(ctx, w)
}

func (r *WeatherForecastTomorrowWidget) ShouldDisplay() bool {
	return true
}

type WeatherPrecipitationWidget struct {
	svc           *services.HomeAssistantService
	precipitation weather.Forecast
}

var _ Widget = (*WeatherPrecipitationWidget)(nil)

func (r *WeatherPrecipitationWidget) Render(ctx context.Context, w io.Writer) {
	m := view.NewWeatherPrecipitationView(r.precipitation)
	WeatherPrecipitation(m).Render(ctx, w)
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
