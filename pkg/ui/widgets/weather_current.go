package widgets

import (
	"context"
	"fmt"

	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type weatherCurrent struct {
	baseWidget
	svc  *services.HomeAssistantService
	font rl.Font
}

func (w *weatherCurrent) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()

	temperature := w.svc.GetWeather().Attributes.Temperature
	rl.DrawTextEx(
		w.font,
		fmt.Sprintf("%d°", temperature),
		rl.NewVector2(0, 0),
		float32(w.font.BaseSize),
		-16.0,
		rl.White,
	)
}

func (w *weatherCurrent) ShouldDisplay() bool {
	return true
}

func NewWeatherCurrent(width, height int32, svc *services.HomeAssistantService) Widget {
	return &weatherCurrent{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:  svc,
		font: rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil, '°'),
	}
}
