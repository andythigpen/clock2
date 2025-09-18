package widgets

import (
	"context"

	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type weatherForecast struct {
	baseWidget
	svc  *services.HomeAssistantService
	font rl.Font
}

func (w *weatherForecast) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()

	rl.DrawTextEx(
		w.font,
		"TODO",
		rl.NewVector2(0, 0),
		float32(w.font.BaseSize),
		-16.0,
		rl.White,
	)
}

func (w *weatherForecast) ShouldDisplay() bool {
	return true
}

func NewWeatherForecast(width, height int32, svc *services.HomeAssistantService) Widget {
	return &weatherForecast{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:  svc,
		font: rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil, '°'),
	}
}
