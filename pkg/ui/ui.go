package ui

import (
	"context"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/platform/screen"
	"github.com/andythigpen/clock2/pkg/services"
	"github.com/andythigpen/clock2/pkg/ui/widgets"
	"github.com/andythigpen/clock2/pkg/ui/widgets/fonts"
)

var (
	frame uint64 = 0
)

func drawWidget(widget widgets.Widget) {
	texture := widget.Texture()
	rl.DrawTexturePro(
		texture,
		rl.NewRectangle(0, 0, float32(texture.Width), -float32(texture.Height)),
		rl.NewRectangle(widget.GetX(), widget.GetY(), float32(texture.Width), float32(texture.Height)),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)
}

func RunForever(haSvc *services.HomeAssistantService, displaySvc *services.DisplayService) {
	rl.InitWindow(platform.ScreenWidth, platform.ScreenHeight, "clock")
	defer rl.CloseWindow()

	rl.SetTargetFPS(platform.FPS)

	texture := rl.LoadRenderTexture(platform.WindowWidth, platform.WindowHeight)

	background := widgets.NewBackground(0, 0, platform.WindowWidth, platform.WindowHeight, haSvc)
	grid := widgets.NewGrid(platform.WindowWidth, platform.WindowHeight)
	clock := widgets.NewClock(0, 0, platform.ClockWidth, platform.WindowHeight)
	carouselWidth := int32(platform.WindowWidth - platform.ClockWidth)
	carouselHeight := int32(platform.WindowHeight)
	carousel := widgets.NewCarousel(
		platform.ClockWidth, 0, carouselWidth, carouselHeight,
		widgets.NewWeatherCurrent(carouselWidth, carouselHeight, haSvc),
		widgets.NewWeatherForecast(carouselWidth, carouselHeight, haSvc),
		widgets.NewWeatherPrecipitation(carouselWidth, carouselHeight, haSvc),
		widgets.NewWeatherTomorrow(carouselWidth, carouselHeight, haSvc),
		widgets.NewHumidity(carouselWidth, carouselHeight, haSvc),
		widgets.NewSun(carouselWidth, carouselHeight, haSvc),
	)
	// ordering matches the render order from back to front
	allWidgets := []widgets.Widget{background, grid, clock, carousel}

	for _, widget := range allWidgets {
		if widget, ok := widget.(widgets.Loader); ok {
			widget.LoadAssets()
		}
	}

	tint := rl.NewColor(255, 255, 255, 255)
	displayState := services.DisplayStateOn

	for !rl.WindowShouldClose() {
		frame += 1

		desiredDisplayState := displaySvc.GetState()
		if desiredDisplayState != displayState {
			switch desiredDisplayState {
			case services.DisplayStateOff:
				screen.ScreenOff()
			case services.DisplayStateOn:
				screen.ScreenOn()
			}
			displayState = desiredDisplayState
		}
		if displayState == services.DisplayStateOff {
			// stop writing to the framebuffer, otherwise raylib will turn the display back on
			time.Sleep(time.Second / platform.FPS)
			continue
		}

		ctx := context.WithValue(context.Background(), widgets.KeyFrame, frame)

		// render widgets to textures first
		for _, w := range allWidgets {
			if f, ok := w.(widgets.Fetcher); ok {
				f.FetchData(ctx)
			}
			if w.ShouldDisplay() {
				w.RenderTexture(ctx)
			}
		}

		// render all textures to a single texture that can be rotated on the actual display
		rl.BeginTextureMode(texture)
		for _, w := range allWidgets {
			if w.ShouldDisplay() {
				drawWidget(w)
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

		desiredBrightness := uint8(rl.Remap(float32(displaySvc.GetBrightness()), 0, 100, 0, 255))
		if tint.R > desiredBrightness {
			tint.R -= 1
			tint.G -= 1
			tint.B -= 1
		} else if tint.R < desiredBrightness {
			tint.R += 1
			tint.G += 1
			tint.B += 1
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.White) // reduces flickering
		rl.DrawTexturePro(texture.Texture, src, dst, rl.NewVector2(0, 0), rotation, tint)
		rl.EndDrawing()
	}

	for _, widget := range allWidgets {
		if widget, ok := widget.(widgets.Loader); ok {
			widget.UnloadAssets()
		}
		widget.Unload()
	}

	fonts.Cache.Unload()

	rl.UnloadRenderTexture(texture)
}
