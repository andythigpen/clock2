package widgets

import (
	"context"
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type weatherForecast struct {
	baseWidget
	svc      *services.HomeAssistantService
	font     rl.Font
	fontHour rl.Font
	icons    map[weather.WeatherCondition]rl.Texture2D
	hours    []forecastHour
}

type forecastHour struct {
	hour        string
	icon        rl.Texture2D
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
				icon:        w.icons[hour.Condition],
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
	if len(w.hours) < 3 {
		panic("expected at least 3 hours of data")
	}

	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()

	width := w.texture.Texture.Width
	iconWidth := w.hours[0].icon.Width
	iconHeight := w.hours[0].icon.Height
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
			rl.NewColor(255, 255, 255, 160),
		)
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
		rl.DrawTexture(w.hours[i].icon, c-iconWidth/2, marginTop, rl.White)
	}
}

func (w *weatherForecast) ShouldDisplay() bool {
	return len(w.hours) >= 3
}

func NewWeatherForecast(width, height int32, svc *services.HomeAssistantService) Widget {
	icons := make(map[weather.WeatherCondition]rl.Texture2D)
	for _, condition := range weather.AllConditions {
		iconName := getWeatherConditionIconName(condition)
		iconPath := getAssetIconPath(iconName, WithSize(256))
		img := rl.LoadImage(iconPath)
		icons[condition] = rl.LoadTextureFromImage(img)
		rl.UnloadImage(img)
	}
	return &weatherForecast{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:      svc,
		font:     rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 240, nil),
		fontHour: rl.LoadFontEx("assets/fonts/Oswald-Bold.ttf", 192, nil),
		icons:    icons,
	}
}
