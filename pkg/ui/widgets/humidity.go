package widgets

import (
	"context"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/services"
)

type humidity struct {
	baseWidget
	svc      *services.HomeAssistantService
	font     rl.Font
	icon     animatedIcon
	humidity string
}

var _ Fetcher = (*humidity)(nil)

func (h *humidity) FetchData(ctx context.Context) {
	currentWeather := h.svc.GetWeather()
	h.humidity = strconv.Itoa(int(currentWeather.Attributes.Humidity))
}

func (h *humidity) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(h.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	// animate the current icon
	x := h.texture.Texture.Width/4 - (h.icon.Width() / 2)
	y := h.texture.Texture.Height/2 - (h.icon.Height() / 2)
	h.icon.RenderFrame(float32(x), float32(y))

	spacing := float32(-16.0)
	textSize := rl.MeasureTextEx(h.font, h.humidity, float32(h.font.BaseSize), spacing)
	textX := float32(h.texture.Texture.Width)*2/3 - textSize.X/2
	rl.DrawTextEx(
		h.font,
		h.humidity,
		rl.NewVector2(textX, -26),
		float32(h.font.BaseSize),
		spacing,
		rl.White,
	)
}

func (h *humidity) ShouldDisplay() bool {
	return true
}

func (h *humidity) LoadAssets() {
	h.icon.LoadAssets()
}

func (h *humidity) UnloadAssets() {
	h.icon.UnloadAssets()
}

func NewHumidity(width, height int32, svc *services.HomeAssistantService) Widget {
	iconPath := getAssetIconPath("humidity", Animated())
	return &humidity{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:  svc,
		font: rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, nil),
		icon: NewAnimatedIcon(iconPath),
	}
}
