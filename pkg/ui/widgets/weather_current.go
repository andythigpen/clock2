package widgets

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"unsafe"

	"github.com/andythigpen/clock2/pkg/models/weather"
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
	prevState    weather.WeatherCondition
	currentState weather.WeatherCondition
	temperature  string
	stateIdx     int
}

var _ Fetcher = (*weatherCurrent)(nil)

func (w *weatherCurrent) FetchData(ctx context.Context) {
	if *uiTestCurrentWeather {
		frame := ctx.Value(KeyFrame).(uint64)
		t := int32(rl.Remap(float32(frame%(30*platform.FPS)), 0, (30*platform.FPS)-1, -20, 110))
		w.temperature = fmt.Sprintf("%d°", t)
		w.currentState = w.prevState
		if frame%360 == 0 { // there are currently 360 frames per animation
			// cycle through icons
			w.stateIdx += 1
			if w.stateIdx > len(weather.AllConditions) {
				w.stateIdx = 0
			}
			w.currentState = weather.AllConditions[w.stateIdx]
		}
	} else {
		currentWeather := w.svc.GetWeather()
		w.currentState = currentWeather.State
		w.temperature = fmt.Sprintf("%d°", currentWeather.Attributes.Temperature)
	}

	// load new icon on state change
	if w.prevState != w.currentState {
		rl.UnloadImage(&w.img)
		w.frameTotal = int32(0)
		w.frameCurrent = 0
		iconName := getWeatherConditionIconName(w.currentState)
		opts := []iconOption{}
		if w.currentState != weather.Exceptional && w.currentState != weather.Unknown {
			opts = append(opts, Animated())
		}
		filename := getAssetIconPath(iconName, opts...)
		w.img = *rl.LoadImageAnim(filename, &w.frameTotal)
		w.icon = rl.LoadTextureFromImage(&w.img)
		w.prevState = w.currentState
	}
}

func (w *weatherCurrent) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	// animate the current icon
	if w.frameTotal > 1 {
		w.frameCurrent += 1
		if w.frameCurrent >= w.frameTotal {
			w.frameCurrent = 0
		}
		dataOffset := w.img.Width * w.img.Height * 4 * w.frameCurrent
		imgSize := w.img.Width * w.img.Height
		rl.UpdateTexture(w.icon,
			unsafe.Slice((*color.RGBA)(unsafe.Pointer(uintptr(w.img.Data)+uintptr(dataOffset))), imgSize))
	}

	rl.DrawTexture(w.icon, 50, 0, rl.White)
	rl.DrawTextEx(
		w.font,
		w.temperature,
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
	iconName := getWeatherConditionIconName(weather.Unknown)
	iconPath := getAssetIconPath(iconName)
	img := rl.LoadImageAnim(iconPath, &frameTotal)
	return &weatherCurrent{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:          svc,
		font:         rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil, '°'),
		img:          *img,
		icon:         rl.LoadTextureFromImage(img),
		frameTotal:   0,
		frameCurrent: 0,
	}
}
