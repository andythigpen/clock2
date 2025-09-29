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
)

var (
	uiTestPrecipitation = flag.Bool("ui-test-precipitation", false, "when set, cycle through weather conditions")
)

type weatherPrecipitation struct {
	baseWidget
	svc           *services.HomeAssistantService
	precipitation *weather.Forecast
	icon          animatedIcon
	font          rl.Font
	fontHour      rl.Font
	prevState     weather.WeatherCondition
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
		idx := int(frame) / 360 % len(weather.AllConditions)
		condition := weather.AllConditions[idx]
		w.precipitation = &weather.Forecast{
			DateTime:                 time.Now(),
			Condition:                condition,
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

	if w.prevState != w.precipitation.Condition {
		iconName := getWeatherConditionIconName(w.precipitation.Condition)
		opts := []iconOption{WithSize(480)}
		if w.precipitation.Condition != weather.Unknown && w.precipitation.Condition != weather.Exceptional {
			opts = append(opts, Animated())
		}
		iconPath := getAssetIconPath(iconName, opts...)
		w.icon.LoadImage(iconPath)
		w.prevState = w.precipitation.Condition
	}
}

func (w *weatherPrecipitation) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	// hour text
	spacing := float32(0)
	textHour := w.precipitation.DateTime.Format("03")
	textSize := rl.MeasureTextEx(w.font, textHour, float32(w.font.BaseSize), spacing)
	textX := float32(w.texture.Texture.Width) / 4
	textY := float32(w.texture.Texture.Height)/4 - textSize.Y/2 + spacing
	rl.DrawTextPro(
		w.fontHour,
		textHour,
		rl.NewVector2(320, 30),
		rl.NewVector2(0, 0),          // origin
		0,                            // rotation
		float32(w.fontHour.BaseSize), // fontSize
		0,                            // spacing
		rl.NewColor(255, 255, 255, 200),
	)

	// animate the current icon
	w.icon.RenderFrame()
	rl.DrawTexture(w.icon.Texture(), 50, 0, rl.White)

	// probability text
	spacing = float32(-16.0)
	textProbability := fmt.Sprintf("%02d", int(w.precipitation.PrecipitationProbability))
	textSize = rl.MeasureTextEx(w.font, textProbability, float32(w.font.BaseSize), spacing)
	textX = float32(w.texture.Texture.Width) / 2
	textY = float32(w.texture.Texture.Height)/2 - textSize.Y/2 + spacing
	rl.DrawTextEx(
		w.font,
		textProbability,
		rl.NewVector2(textX, textY),
		float32(w.font.BaseSize),
		spacing,
		rl.White,
	)

}

func (w *weatherPrecipitation) ShouldDisplay() bool {
	return w.precipitation != nil
}

func NewWeatherPrecipitation(width, height int32, svc *services.HomeAssistantService) Widget {
	return &weatherPrecipitation{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:      svc,
		font:     rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil),
		fontHour: rl.LoadFontEx("assets/fonts/Oswald-Bold.ttf", 192, nil),
	}
}
