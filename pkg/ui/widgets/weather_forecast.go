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

type weatherForecast struct {
	baseWidget
	svc      *services.HomeAssistantService
	font     rl.Font
	fontHour rl.Font
	icons    []rl.Texture2D
	hours    []forecastHour
}

type forecastHour struct {
	hour        string
	hour24      int
	condition   weather.WeatherCondition
	temperature string
}

func (w *weatherForecast) FetchData(ctx context.Context) {
	now := time.Now().Local()
	hours := []forecastHour{}
	forecast := w.svc.GetForecast()
	for _, hour := range forecast.Attributes.Forecast {
		if hour.DateTime.After(now) {
			hours = append(hours, forecastHour{
				hour:        hour.DateTime.Format("03"),
				hour24:      hour.DateTime.Hour(),
				condition:   hour.Condition,
				temperature: strconv.Itoa(int(hour.Temperature)),
			})
			if len(hours) >= 3 {
				break
			}
		}
	}
	w.hours = hours
}

func (w *weatherForecast) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	if len(w.icons) == 0 || len(w.hours) == 0 {
		// data may not be loaded yet...nothing to draw
		return
	}

	width := w.texture.Texture.Width
	margin := int32(platform.Margin)
	marginTop := int32(35)
	marginTopText := int32(-45)

	centers := []int32{
		(width / 6) + margin/2,
		(width / 2),
		(width * 5 / 6) - margin/2,
	}

	for i, c := range centers {
		textSize := rl.MeasureTextEx(w.font, w.hours[i].temperature, float32(w.font.BaseSize), 0)
		rl.DrawTextPro(
			w.fontHour,
			w.hours[i].hour,
			rl.NewVector2(float32(c), 0),
			rl.NewVector2(0, 0),          // origin
			0,                            // rotation
			float32(w.fontHour.BaseSize), // fontSize
			0,                            // spacing
			rl.NewColor(255, 255, 255, 200),
		)
		iconHeight := w.icons[i].Height
		rl.DrawTextPro(
			w.font,
			w.hours[i].temperature,
			rl.NewVector2(float32(c)-(textSize.X/2), float32(marginTopText+iconHeight)),
			rl.NewVector2(0, 0),      // origin
			0,                        // rotation
			float32(w.font.BaseSize), // fontSize
			0,                        // spacing
			rl.White,
		)
		iconWidth := w.icons[i].Width
		rl.DrawTexture(w.icons[i], c-iconWidth/2, marginTop, rl.White)
	}
}

func (w *weatherForecast) ShouldDisplay() bool {
	return len(w.hours) >= 3
}

func (w *weatherForecast) LoadAssets() {
	for _, h := range w.hours {
		iconType := icons.GetWeatherConditionIconType(h.condition, icons.WithHourOfDay(h.hour24))
		iconPath := icons.GetStaticIconPath(iconType, icons.WithSize(256))
		icon := rl.LoadTexture(iconPath)
		w.icons = append(w.icons, icon)
	}
}

func (w *weatherForecast) UnloadAssets() {
	for _, i := range w.icons {
		rl.UnloadTexture(i)
	}
	w.icons = nil
}

func NewWeatherForecast(width, height int32, svc *services.HomeAssistantService) Widget {
	return &weatherForecast{
		baseWidget: newBaseWidget(0, 0, width, height),
		svc:        svc,
		font:       fonts.Cache.Load(fonts.FontOswald, 240),
		fontHour:   fonts.Cache.Load(fonts.FontOswald, 192, fonts.WithVariation(fonts.FontVariationBold)),
	}
}
