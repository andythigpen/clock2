package widgets

import (
	"context"
	"flag"
	"fmt"

	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	uiTestCurrentWeather = flag.Bool("ui-test-current-weather", false, "when true, cycle through current states")
)

type weatherCurrent struct {
	baseWidget
	svc  *services.HomeAssistantService
	font rl.Font
	icon rl.Texture2D
}

func (w *weatherCurrent) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	var (
		temperature string
		icon        rl.Texture2D
	)

	if *uiTestCurrentWeather {
		frame := ctx.Value(KeyFrame).(uint64)
		t := int32(rl.Remap(float32(frame%(30*platform.FPS)), 0, (30*platform.FPS)-1, -20, 110))
		temperature = fmt.Sprintf("%d°", t)
		// TODO cycle through icons
		icon = w.icon
	} else {
		temperature = fmt.Sprintf("%d°", w.svc.GetWeather().Attributes.Temperature)
		// TODO select icon based on state
		icon = w.icon
	}

	rl.DrawTexture(icon, 50, 0, rl.White)
	rl.DrawTextEx(
		w.font,
		temperature,
		rl.NewVector2(540, -26),
		float32(w.font.BaseSize),
		-16.0,
		rl.White,
	)
}

func (w *weatherCurrent) ShouldDisplay() bool {
	return true
}

func NewWeatherCurrent(width, height int32, svc *services.HomeAssistantService) Widget {
	clearDay := rl.LoadImage("assets/icons/weather/static/png/480/clear-day.png")
	return &weatherCurrent{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:  svc,
		font: rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil, '°'),
		icon: rl.LoadTextureFromImage(clearDay),
	}
}
