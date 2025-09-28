package widgets

import (
	"context"
	"strconv"
	"time"

	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type weatherForecast struct {
	baseWidget
	svc      *services.HomeAssistantService
	font     rl.Font
	fontHour rl.Font
	icon     rl.Texture2D
}

type forecastHour struct {
	hour        string
	icon        rl.Texture2D
	temperature string
}

func (w *weatherForecast) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(w.texture)
	defer rl.EndTextureMode()

	now := time.Now()
	hours := []forecastHour{}
	forecast := w.svc.GetForecast()
	for _, hour := range forecast.Attributes.Forecast {
		if hour.DateTime.After(now) {
			hours = append(hours, forecastHour{
				hour:        hour.DateTime.Format("03"),
				icon:        w.icon, // TODO
				temperature: strconv.Itoa(int(hour.Temperature)),
			})
			if len(hours) >= 3 {
				break
			}
		}
	}
	if len(hours) < 3 {
		// TODO
		return
	}

	width := w.texture.Texture.Width
	iconWidth := hours[0].icon.Width
	iconHeight := hours[0].icon.Height
	margin := int32(platform.Margin)
	marginTop := int32(35)
	marginTopText := int32(-45)

	centers := []int32{
		(width / 6) + margin/2,
		(width / 2),
		(width * 5 / 6) - margin/2,
	}

	for i, c := range centers {
		textSize := rl.MeasureTextEx(w.font, hours[i].temperature, float32(w.font.BaseSize), 0)
		rl.DrawTextPro(
			w.fontHour,
			hours[i].hour,
			rl.NewVector2(float32(c), 0),
			rl.NewVector2(0, 0),          // origin
			0,                            // rotation
			float32(w.fontHour.BaseSize), // fontSize
			0,                            // spacing
			rl.NewColor(255, 255, 255, 160),
		)
		rl.DrawTextPro(
			w.font,
			hours[i].temperature,
			rl.NewVector2(float32(c)-(textSize.X/2), float32(marginTopText+iconHeight)),
			rl.NewVector2(0, 0),      // origin
			0,                        // rotation
			float32(w.font.BaseSize), // fontSize
			0,                        // spacing
			rl.White,
		)
		rl.DrawTexture(hours[i].icon, c-iconWidth/2, marginTop, rl.White)
	}

	// rl.DrawTexture(hours[0].icon, (width/6)-iconWidth/2+margin/2, marginTop, rl.White)
	// rl.DrawTexture(hours[1].icon, (width/2)-iconWidth/2, marginTop, rl.White)
	// rl.DrawTexture(hours[2].icon, (width*5/6)-iconWidth/2-margin/2, marginTop, rl.White)

	// rl.MeasureTextEx(w.font)
	// rl.DrawText()
	// rl.DrawTextEx(
	// 	w.font,
	// 	"TODO",
	// 	rl.NewVector2(0, 0),
	// 	float32(w.font.BaseSize),
	// 	-16.0,
	// 	rl.White,
	// )
}

func (w *weatherForecast) ShouldDisplay() bool {
	return false
}

func NewWeatherForecast(width, height int32, svc *services.HomeAssistantService) Widget {
	clearDay := rl.LoadImage("assets/icons/weather/static/png/256/clear-day.png")
	return &weatherForecast{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:      svc,
		font:     rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 240, nil),
		fontHour: rl.LoadFontEx("assets/fonts/Oswald-Bold.ttf", 192, nil),
		icon:     rl.LoadTextureFromImage(clearDay),
	}
}
