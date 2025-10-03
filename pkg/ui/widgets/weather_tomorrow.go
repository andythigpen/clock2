package widgets

import (
	"context"
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/services"
	"github.com/andythigpen/clock2/pkg/ui/widgets/fonts"
	"github.com/andythigpen/clock2/pkg/ui/widgets/icons"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type weatherTomorrow struct {
	baseWidget
	svc          *services.HomeAssistantService
	prevState    weather.WeatherCondition
	currentState weather.WeatherCondition
	icon         icons.AnimatedIcon
	font         rl.Font
	fontTitle    rl.Font
	tempLo       int8
	tempHi       int8
}

var _ Fetcher = (*weatherTomorrow)(nil)

var gradients = []rl.Color{
	rl.NewColor(36, 27, 99, 255),    // "#241b63" <30
	rl.NewColor(32, 52, 141, 255),   // "#20348d" 30
	rl.NewColor(34, 67, 150, 255),   // "#224396" 35
	rl.NewColor(61, 100, 174, 255),  // "#3d64ae" 40
	rl.NewColor(85, 131, 193, 255),  // "#5583c1" 45
	rl.NewColor(137, 173, 220, 255), // "#89addc" 50
	rl.NewColor(155, 186, 226, 255), // "#9bbae2" 55
	rl.NewColor(154, 205, 207, 255), // "#9acdcf" 60
	rl.NewColor(162, 207, 164, 255), // "#a2cfa4" 65
	rl.NewColor(215, 222, 126, 255), // "#d7de7e" 70
	rl.NewColor(244, 216, 98, 255),  // "#f4d862" 75
	rl.NewColor(242, 148, 0, 255),   // "#f29400" 80
	rl.NewColor(235, 86, 30, 255),   // "#eb561e" 85
	rl.NewColor(195, 5, 7, 255),     // "#c30507" 90
	rl.NewColor(112, 3, 24, 255),    // "#700318" >90
}

func getGradientColor(temp int8) rl.Color {
	if temp < 30 {
		return gradients[0]
	}
	if temp > 90 {
		return gradients[len(gradients)-1]
	}
	idx := int32(rl.Remap(float32(temp), 30, 90, 1, float32(len(gradients)-2)))
	return gradients[idx]
}

func (w *weatherTomorrow) FetchData(ctx context.Context) {
	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	dayAfterTomorrow := time.Now().AddDate(0, 0, 2).Truncate(24 * time.Hour)
	conditions := map[weather.WeatherCondition]int{}
	w.tempHi = int8(-127)
	w.tempLo = int8(127)
	forecast := w.svc.GetForecast()
	for _, hour := range forecast.Attributes.Forecast {
		if !hour.DateTime.After(tomorrow) {
			continue
		}
		if !hour.DateTime.Before(dayAfterTomorrow) {
			break
		}
		conditions[hour.Condition] += 1
		if hour.Temperature < w.tempLo {
			w.tempLo = hour.Temperature
		} else if hour.Temperature > w.tempHi {
			w.tempHi = hour.Temperature
		}
	}
	w.currentState = weather.Unknown
	maximum := 0
	for c, num := range conditions {
		if num > maximum {
			w.currentState = c
			maximum = num
		}
	}
	if w.prevState != w.currentState {
		iconType := icons.GetWeatherConditionIconType(w.currentState)
		w.icon.SetIconType(iconType)
		w.prevState = w.currentState
	}
}

func (w *weatherTomorrow) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	// animate the current icon
	x := w.texture.Texture.Width/4 - (w.icon.Width() / 2)
	y := w.texture.Texture.Height/2 - (w.icon.Height() / 2)
	w.icon.RenderFrame(float32(x), float32(y))

	spacing := float32(-6.0)
	title := "Tomorrow"
	titleSize := rl.MeasureTextEx(w.fontTitle, title, float32(w.fontTitle.BaseSize), spacing)
	titleX := float32(w.texture.Texture.Width)*3/4 - (titleSize.X / 2) - platform.Margin/2
	tempLo := strconv.Itoa(int(w.tempLo))
	tempLoSize := rl.MeasureTextEx(w.font, tempLo, float32(w.font.BaseSize), spacing)
	tempLoX := float32(w.texture.Texture.Width)/2 + spacing
	tempHi := strconv.Itoa(int(w.tempHi))
	tempHiSize := rl.MeasureTextEx(w.font, tempHi, float32(w.font.BaseSize), spacing)
	tempHiX := float32(w.texture.Texture.Width) - platform.Margin - tempHiSize.X
	tempY := float32(w.texture.Texture.Height) - 65 - tempLoSize.Y
	rl.DrawTextEx(
		w.fontTitle,
		title,
		rl.NewVector2(titleX, 50),     // position
		float32(w.fontTitle.BaseSize), // font-size
		spacing,                       // spacing
		rl.White,
	)
	rl.DrawTextEx(
		w.font,
		strconv.Itoa(int(w.tempLo)),
		rl.NewVector2(tempLoX, tempY), // position
		float32(w.font.BaseSize),      // font-size
		spacing,                       // spacing
		rl.White,
	)
	rl.DrawTextEx(
		w.font,
		strconv.Itoa(int(w.tempHi)),
		rl.NewVector2(tempHiX, tempY), // position
		float32(w.font.BaseSize),      // font-size
		spacing,                       // spacing
		rl.White,
	)
	colorLo := getGradientColor(w.tempLo)
	colorHi := getGradientColor(w.tempHi)
	colorHeight := float32(30)
	colorX := tempLoX + tempLoSize.X + platform.Margin
	colorY := tempY + tempLoSize.Y/2 - colorHeight/2
	colorWidth := w.texture.Texture.Width/2 - int32(tempLoSize.X) - int32(tempHiSize.X) - platform.Margin*2 - 10
	rl.DrawRectangleGradientEx(
		rl.NewRectangle(colorX, colorY, float32(colorWidth), colorHeight),
		colorLo, // top left
		colorLo, // bottom left
		colorHi, // top right
		colorHi, // bottom right
	)
	radius := colorHeight / 2
	rl.DrawCircleSector(rl.NewVector2(colorX, colorY+radius), radius, 90, 270, 10, colorLo)
	rl.DrawCircleSector(rl.NewVector2(colorX+float32(colorWidth), colorY+radius), radius, -90, 90, 10, colorHi)
}

func (w *weatherTomorrow) ShouldDisplay() bool {
	return time.Now().Hour() >= 17 && w.currentState != weather.Unknown
}

func (w *weatherTomorrow) LoadAssets() {
	w.icon.LoadAssets()
}

func (w *weatherTomorrow) UnloadAssets() {
	w.icon.UnloadAssets()
}

func NewWeatherTomorrow(width, height int32, svc *services.HomeAssistantService) Widget {
	iconType := icons.GetWeatherConditionIconType(weather.Unknown)
	return &weatherTomorrow{
		baseWidget: newBaseWidget(0, 0, width, height),
		svc:        svc,
		font:       fonts.Cache.Load(fonts.FontOswald, 240),
		fontTitle:  fonts.Cache.Load(fonts.FontOswald, 192, fonts.WithVariation(fonts.FontVariationBold)),
		icon:       icons.NewAnimatedIcon(iconType),
	}
}
