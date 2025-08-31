package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/andythigpen/clock2/pkg/components"
	. "github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/services"
)

func Register(mux *http.ServeMux, haSvc *services.HomeAssistantService) {
	component := components.Index()
	mux.Handle("/", templ.Handler(component))
	mux.Handle("/components/weather-current", &WeatherCurrentHandler{svc: haSvc})
	mux.Handle("/components/weather-forecast", &WeatherForecastHandler{svc: haSvc})
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
}

type WeatherCurrentHandler struct {
	svc *services.HomeAssistantService
}

func (h *WeatherCurrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	weather := h.svc.GetWeather()
	forecast := h.svc.GetForecast()
	m := view.NewWeatherCurrentView(weather, forecast)
	WeatherCurrent(m).Render(context.Background(), w)
}

type WeatherForecastHandler struct {
	svc *services.HomeAssistantService
}

func (h *WeatherForecastHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	forecast := h.svc.GetForecast()
	m := view.NewWeatherForecastView(forecast)
	WeatherForecast(m).Render(context.Background(), w)
}
