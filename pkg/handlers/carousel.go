package handlers

import (
	"context"
	"net/http"

	"github.com/andythigpen/clock2/pkg/handlers/widgets"
	"github.com/andythigpen/clock2/pkg/services"
)

type CarouselHandler struct {
	widgets []widgets.Widget
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
		widgets: []widgets.Widget{
			widgets.NewWeatherCurrentWidget(svc),
			widgets.NewWeatherForecastWidget(svc),
			widgets.NewWeatherPrecipitationWidget(svc),
			widgets.NewWeatherHumidityWidget(svc),
			widgets.NewSunWidget(svc),
			widgets.NewWeatherForecastTomorrowWidget(svc),
		},
	}
}
