package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/andythigpen/clock2/pkg/components"
	"github.com/andythigpen/clock2/pkg/services"
)

func Register(mux *http.ServeMux, haSvc *services.HomeAssistantService) {
	component := components.Index()
	mux.Handle("/", templ.Handler(component))
	mux.Handle("/carousel", &CarouselHandler{
		widgets: []Widget{
			&WeatherCurrentWidget{svc: haSvc},
			&WeatherForecastWidget{svc: haSvc},
			&WeatherForecastTomorrowWidget{svc: haSvc},
			// TODO: humidity
			// TODO: precipitation
			// TODO: sun
		},
	})
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
}

type CarouselHandler struct {
	widgets []Widget
	cursor  int
}

func (h *CarouselHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	for range len(h.widgets) {
		idx := h.cursor
		h.cursor += 1
		if h.cursor >= len(h.widgets) {
			h.cursor = 0
		}
		widget := h.widgets[idx]
		if widget.ShouldDisplay() {
			widget.Render(ctx, w)
			break
		}
	}
}
