package widgets

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"unsafe"

	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	uiTestCurrentWeather = flag.Bool("ui-test-current-weather", false, "when true, cycle through current states")
)

type weatherCurrent struct {
	baseWidget
	svc          *services.HomeAssistantService
	font         rl.Font
	img          rl.Image
	icon         rl.Texture2D
	frameCurrent int32
	frameTotal   int32
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

	w.frameCurrent += 1
	if w.frameCurrent >= w.frameTotal {
		w.frameCurrent = 0
	}
	dataOffset := w.img.Width * w.img.Height * 4 * w.frameCurrent
	imgSize := w.img.Width * w.img.Height
	rl.UpdateTexture(w.icon,
		unsafe.Slice((*color.RGBA)(unsafe.Pointer(uintptr(w.img.Data)+uintptr(dataOffset))), imgSize))

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
	frameTotal := int32(0)
	img := rl.LoadImageAnim("assets/icons/weather/animated/gif/480/thunderstorms-day.gif", &frameTotal)
	return &weatherCurrent{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:          svc,
		font:         rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil, '°'),
		img:          *img,
		icon:         rl.LoadTextureFromImage(img),
		frameTotal:   frameTotal,
		frameCurrent: 0,
	}
}
