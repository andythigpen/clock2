package ui

import (
	"context"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/services"
	"github.com/andythigpen/clock2/pkg/ui/widgets"
)

var (
	frame uint64 = 0
)

func drawWidget(widget widgets.Widget, a uint8) {
	texture := widget.Texture()
	rl.DrawTexturePro(
		texture,
		rl.NewRectangle(0, 0, float32(texture.Width), -float32(texture.Height)),
		rl.NewRectangle(widget.GetX(), widget.GetY(), float32(texture.Width), float32(texture.Height)),
		rl.NewVector2(0, 0),
		0,
		rl.NewColor(255, 255, 255, a),
	)
}

func RunForever(haSvc *services.HomeAssistantService) {
	rl.InitWindow(platform.ScreenWidth, platform.ScreenHeight, "clock")
	defer rl.CloseWindow()

	rl.SetTargetFPS(platform.FPS)

	texture := rl.LoadRenderTexture(platform.WindowWidth, platform.WindowHeight)

	background := widgets.NewBackground(0, 0, platform.WindowWidth, platform.WindowHeight)
	grid := widgets.NewGrid(platform.WindowWidth, platform.WindowHeight)
	clock := widgets.NewClock(0, 0, platform.ClockWidth, platform.WindowHeight)
	carouselWidth := int32(platform.WindowWidth - platform.ClockWidth)
	carousel := widgets.NewCarousel(platform.ClockWidth, 0,
		widgets.NewWeatherCurrent(carouselWidth, platform.WindowHeight, haSvc),
		widgets.NewWeatherForecast(carouselWidth, platform.WindowHeight, haSvc),
	)
	// ordering matches the render order from back to front
	allWidgets := []widgets.Widget{background, grid, clock, carousel}

	for !rl.WindowShouldClose() {
		frame += 1
		ctx := context.WithValue(context.Background(), widgets.KeyFrame, frame)

		// render widgets to textures first
		for _, w := range allWidgets {
			if w.ShouldDisplay() {
				w.RenderTexture(ctx)
			}
		}

		// render all textures to a single texture that can be rotated on the actual display
		rl.BeginTextureMode(texture)
		for idx, w := range allWidgets {
			if w.ShouldDisplay() {
				a := uint8(255)
				if idx == 3 {
					a = uint8(rl.Remap(float32(frame%300), 0, 600, 0, 255))
				}
				drawWidget(w, a)
			}
		}
		rl.EndTextureMode()

		var (
			src      rl.Rectangle
			dst      rl.Rectangle
			rotation float32
		)
		switch platform.Platform {
		case platform.PlatformDesktop:
			src = rl.NewRectangle(0, 0, float32(texture.Texture.Width), -float32(texture.Texture.Height))
			dst = rl.NewRectangle(0, 0, float32(texture.Texture.Width), float32(texture.Texture.Height))
			rotation = 0.0
		case platform.PlatformPI:
			src = rl.NewRectangle(0, 0, -float32(texture.Texture.Width), float32(texture.Texture.Height))
			dst = rl.NewRectangle(float32(texture.Texture.Height), 0, float32(texture.Texture.Width), float32(texture.Texture.Height))
			rotation = 90.0
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.White) // reduces flickering
		rl.DrawTexturePro(texture.Texture, src, dst, rl.NewVector2(0, 0), rotation, rl.White)
		rl.EndDrawing()
	}
}
