package widgets

import (
	"context"

	"github.com/andythigpen/clock2/pkg/services"
)

var _ Widget = (*WeatherCurrentWidget)(nil)

type WeatherCurrentWidget struct {
	svc *services.HomeAssistantService
}

func (w *WeatherCurrentWidget) Initialize() {
	panic("unimplemented")
}

func (w *WeatherCurrentWidget) Render(ctx context.Context, x, y, width, height float32) {
	panic("unimplemented")
}

func (w *WeatherCurrentWidget) ShouldDisplay() bool {
	panic("unimplemented")
}

func NewWeatherCurrentWidget(svc *services.HomeAssistantService) Widget {
	return &WeatherCurrentWidget{svc}
}
