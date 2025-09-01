package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/andythigpen/clock2/pkg/handlers/widgets"
	"github.com/andythigpen/clock2/pkg/services"
)

type CarouselHandler struct {
	displaySvc *services.DisplayService
	widgets    []widgets.Widget
	cursor     int
}

func (h *CarouselHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state, err := h.displaySvc.GetState()
	if err != nil {
		slog.Error("failed to get display state", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if state == services.DisplayStateOff {
		// don't show anything when the display is off
		w.WriteHeader(http.StatusOK)
		return
	}

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

func NewCarouselHandler(haSvc *services.HomeAssistantService, displaySvc *services.DisplayService) *CarouselHandler {
	return &CarouselHandler{
		displaySvc: displaySvc,
		widgets: []widgets.Widget{
			widgets.NewWeatherCurrentWidget(haSvc),
			widgets.NewWeatherForecastWidget(haSvc),
			widgets.NewWeatherPrecipitationWidget(haSvc),
			widgets.NewWeatherHumidityWidget(haSvc),
			widgets.NewSunWidget(haSvc),
			widgets.NewWeatherForecastTomorrowWidget(haSvc),
		},
	}
}
