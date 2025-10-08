package widgets

import (
	"context"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/services"
	"github.com/andythigpen/clock2/pkg/ui/widgets/fonts"
	"github.com/andythigpen/clock2/pkg/ui/widgets/icons"
)

type humidity struct {
	baseWidget
	svc         *services.HomeAssistantService
	font        rl.Font
	fontPercent rl.Font
	icon        icons.AnimatedIcon
	humidity    string
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
	textX := float32(h.texture.Texture.Width)/2 + spacing
	rl.DrawTextEx(
		h.font,
		h.humidity,
		rl.NewVector2(textX, -26),
		float32(h.font.BaseSize),
		spacing,
		rl.White,
	)
	textSize := rl.MeasureTextEx(h.font, h.humidity, float32(h.font.BaseSize), spacing)
	rl.DrawTextEx(
		h.fontPercent,
		"%",
		rl.NewVector2(textX+textSize.X, 50),
		float32(h.fontPercent.BaseSize),
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

func (h *humidity) Unload() {
	h.baseWidget.Unload()
	rl.UnloadFont(h.font)
}

func NewHumidity(width, height int32, svc *services.HomeAssistantService) Widget {
	return &humidity{
		baseWidget:  newBaseWidget(0, 0, width, height),
		svc:         svc,
		font:        rl.LoadFontEx(fonts.GetAssetFontPath(fonts.FontOswald), 500, fonts.Numbers),
		fontPercent: fonts.Cache.Load(fonts.FontOswald, 192),
		icon:        icons.NewAnimatedIcon(icons.IconHumidity),
	}
}
