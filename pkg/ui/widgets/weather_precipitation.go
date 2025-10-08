package widgets

import (
	"context"
	"flag"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/models/weather"
	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/services"
	"github.com/andythigpen/clock2/pkg/ui/widgets/fonts"
	"github.com/andythigpen/clock2/pkg/ui/widgets/icons"
)

var (
	uiTestPrecipitation = flag.Bool("ui-test-precipitation", false, "when set, cycle through weather conditions")
)

type weatherPrecipitation struct {
	baseWidget
	svc                    *services.HomeAssistantService
	precipitation          *weather.Forecast
	testPrecipitationIndex int
	icon                   icons.AnimatedIcon
	font                   rl.Font
	fontHour               rl.Font
	fontPercent            rl.Font
}

var _ Fetcher = (*weatherPrecipitation)(nil)

func isPrecipitation(condition weather.WeatherCondition) bool {
	switch condition {
	case weather.Rain, weather.Thunderstorms, weather.ThunderstormsRain, weather.Sleet, weather.Snow, weather.Hail:
		return true
	default:
		return false
	}
}

func (w *weatherPrecipitation) FetchData(ctx context.Context) {
	if *uiTestPrecipitation {
		frame := ctx.Value(KeyFrame).(uint64)
		w.precipitation = &weather.Forecast{
			DateTime:                 time.Now().Add(time.Duration(frame/platform.FPS%24) * time.Hour),
			PrecipitationProbability: uint8(frame / platform.FPS % 100),
		}
	} else {
		forecast := w.svc.GetForecast()
		count := 0
		for _, hour := range forecast.Attributes.Forecast {
			if hour.DateTime.Local().Before(time.Now()) {
				continue
			}
			count += 1
			if count >= 4 {
				w.precipitation = nil
				return
			}
			if !isPrecipitation(hour.Condition) || hour.PrecipitationProbability < 10 {
				continue
			}
			w.precipitation = &hour
			break
		}
	}

	if w.precipitation == nil {
		return
	}

	iconType := icons.GetWeatherConditionIconType(w.precipitation.Condition, icons.WithHourOfDay(w.precipitation.DateTime.Hour()))
	w.icon.SetIconType(iconType)
}

func (w *weatherPrecipitation) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	// hour text
	spacing := float32(0)
	textHour := w.precipitation.DateTime.Format("03")
	rl.DrawTextPro(
		w.fontHour,
		textHour,
		rl.NewVector2(350, 0),
		rl.NewVector2(0, 0),          // origin
		0,                            // rotation
		float32(w.fontHour.BaseSize), // fontSize
		spacing,                      // spacing
		rl.NewColor(255, 255, 255, 200),
	)

	// animate the current icon
	x := w.texture.Texture.Width/4 - (w.icon.Width() / 2)
	y := w.texture.Texture.Height/2 - (w.icon.Height() / 2)
	w.icon.RenderFrame(float32(x), float32(y))

	// probability text
	spacing = float32(-16.0)
	textProbability := fmt.Sprintf("%02d", int(w.precipitation.PrecipitationProbability))
	textSize := rl.MeasureTextEx(w.font, textProbability, float32(w.font.BaseSize), spacing)
	textX := float32(w.texture.Texture.Width)/2 + spacing
	textY := float32(w.texture.Texture.Height)/2 - textSize.Y/2 + spacing
	rl.DrawTextEx(
		w.font,
		textProbability,
		rl.NewVector2(textX, textY),
		float32(w.font.BaseSize),
		spacing,
		rl.White,
	)
	rl.DrawTextEx(
		w.fontPercent,
		"%",
		rl.NewVector2(textX+textSize.X, 50),
		float32(w.fontPercent.BaseSize),
		spacing,
		rl.White,
	)
}

func (w *weatherPrecipitation) ShouldDisplay() bool {
	return w.precipitation != nil
}

func (w *weatherPrecipitation) LoadAssets() {
	if *uiTestPrecipitation {
		condition := weather.PrecipitationConditions[w.testPrecipitationIndex]
		iconType := icons.GetWeatherConditionIconType(condition)
		w.icon.SetIconType(iconType)
		w.testPrecipitationIndex += 1
		if w.testPrecipitationIndex >= len(weather.PrecipitationConditions) {
			w.testPrecipitationIndex = 0
		}
	}
	w.icon.LoadAssets()
}

func (w *weatherPrecipitation) UnloadAssets() {
	w.icon.UnloadAssets()
}

func (w *weatherPrecipitation) Unload() {
	w.baseWidget.Unload()
	rl.UnloadFont(w.font)
}

func NewWeatherPrecipitation(width, height int32, svc *services.HomeAssistantService) Widget {
	iconType := icons.GetWeatherConditionIconType(weather.Unknown)
	return &weatherPrecipitation{
		baseWidget:  newBaseWidget(0, 0, width, height),
		svc:         svc,
		font:        rl.LoadFontEx(fonts.GetAssetFontPath(fonts.FontOswald), 500, fonts.Numbers),
		fontHour:    fonts.Cache.Load(fonts.FontOswald, 192, fonts.WithVariation(fonts.FontVariationBold)),
		fontPercent: fonts.Cache.Load(fonts.FontOswald, 192),
		icon:        icons.NewAnimatedIcon(iconType),
	}
}
