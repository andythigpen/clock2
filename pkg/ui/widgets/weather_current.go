package widgets

import (
	"context"
	"flag"
	"fmt"
	"time"

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
	svc                  *services.HomeAssistantService
	font                 rl.Font
	icon                 animatedIcon
	iconRising           animatedIcon
	iconFalling          animatedIcon
	prevState            weather.WeatherCondition
	currentState         weather.WeatherCondition
	temperature          string
	temperatureDirection int
	stateIdx             int
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
		currentTemp := currentWeather.Attributes.Temperature
		w.temperature = fmt.Sprintf("%d°", currentTemp)
		w.temperatureDirection = 0
		forecast := w.svc.GetForecast()
		now := time.Now()
		temps := []int8{}
		for _, hour := range forecast.Attributes.Forecast {
			if hour.DateTime.Before(now) {
				continue
			}
			temps = append(temps, hour.Temperature)
			if len(temps) >= 2 {
				break
			}
		}
		if len(temps) >= 2 {
			w.temperatureDirection = int(temps[1] - temps[0])
			if w.temperatureDirection < 0 && currentTemp-temps[0] <= 2 {
				w.temperatureDirection = 0
			} else if w.temperatureDirection > 0 && temps[0]-currentTemp <= 2 {
				w.temperatureDirection = 0
			}
		}
	}

	// load new icon on state change
	if w.prevState != w.currentState {
		iconName := getWeatherConditionIconName(w.currentState)
		filename := getAssetIconPath(iconName, Animated())
		w.icon.SetFilename(filename)
		w.prevState = w.currentState
	}
}

func (w *weatherCurrent) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	// animate the current icon
	x := w.texture.Texture.Width/4 - (w.icon.Width() / 2)
	y := w.texture.Texture.Height/2 - (w.icon.Height() / 2)
	w.icon.RenderFrame(float32(x), float32(y))

	spacing := float32(-16)
	textSize := rl.MeasureTextEx(w.font, w.temperature, float32(w.font.BaseSize), spacing)
	textX := float32(w.texture.Texture.Width)/2 + spacing
	textY := float32(-26)

	// animate the direction arrow
	if w.temperatureDirection != 0 {
		var icon *animatedIcon
		if w.temperatureDirection < 0 {
			icon = &w.iconFalling
		} else {
			icon = &w.iconRising
		}
		degreeSize := rl.MeasureTextEx(w.font, "°", float32(w.font.BaseSize), spacing)
		directionX := textX + textSize.X - degreeSize.X + (degreeSize.X-float32(icon.Width()))/2 + 1
		directionY := textY + 230 // degreeSize.Y is the entire line size, which we can't use here
		icon.RenderFrame(directionX, directionY)
	}

	rl.DrawTextEx(
		w.font,
		w.temperature,
		rl.NewVector2(textX, textY),
		float32(w.font.BaseSize),
		spacing,
		rl.White,
	)
}

func (w *weatherCurrent) ShouldDisplay() bool {
	return true
}

func (w *weatherCurrent) LoadAssets() {
	w.icon.LoadAssets()
	w.iconRising.LoadAssets()
	w.iconFalling.LoadAssets()
}

func (w *weatherCurrent) UnloadAssets() {
	w.icon.UnloadAssets()
	w.iconRising.UnloadAssets()
	w.iconFalling.UnloadAssets()
}

func (w *weatherCurrent) Unload() {
	w.baseWidget.Unload()
	rl.UnloadFont(w.font)
}

func NewWeatherCurrent(width, height int32, svc *services.HomeAssistantService) Widget {
	iconName := getWeatherConditionIconName(weather.Unknown)
	iconPath := getAssetIconPath(iconName, Animated())
	return &weatherCurrent{
		baseWidget: newBaseWidget(0, 0, width, height),
		svc:        svc,
		// not from the font cache because of the extra rune
		font:        rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil, '°'),
		icon:        NewAnimatedIcon(iconPath),
		iconRising:  NewAnimatedIcon(getAssetIconPath("pressure-high", Animated())),
		iconFalling: NewAnimatedIcon(getAssetIconPath("pressure-low", Animated())),
	}
}
