package widgets

import (
	"context"
	"io"
	"time"

	"github.com/andythigpen/clock2/pkg/components/weather"
	"github.com/andythigpen/clock2/pkg/models/view"
	"github.com/andythigpen/clock2/pkg/services"
)

type SunWidget struct {
	svc *services.HomeAssistantService
}

func NewSunWidget(svc *services.HomeAssistantService) Widget {
	return &SunWidget{svc}
}

var _ Widget = (*SunWidget)(nil)

func (r *SunWidget) Render(ctx context.Context, w io.Writer) {
	sun := r.svc.GetSun()
	m := view.NewSunView(sun)
	weather.Sun(m).Render(ctx, w)
}

func (r *SunWidget) ShouldDisplay() bool {
	sun := r.svc.GetSun()
	if sun.Attributes.NextRising.Before(time.Now().Add(2 * time.Hour)) {
		return true
	} else if sun.Attributes.NextSetting.Before(time.Now().Add(2 * time.Hour)) {
		return true
	}
	return false
}
