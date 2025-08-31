package handlers

import (
	"context"
	"net/http"

	"github.com/andythigpen/clock2/pkg/services"
)

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

func NewCarouselHandler(svc *services.HomeAssistantService) *CarouselHandler {
	return &CarouselHandler{
		widgets: []Widget{
			&WeatherCurrentWidget{svc: svc},
			&WeatherForecastWidget{svc: svc},
			&WeatherPrecipitationWidget{svc: svc},
			&WeatherHumidityWidget{svc: svc},
			&SunWidget{svc: svc},
			&WeatherForecastTomorrowWidget{svc: svc},
		},
	}
}
