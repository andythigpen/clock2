package ui

import (
	"context"
	"flag"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/platform"
	"github.com/andythigpen/clock2/pkg/ui/widgets"
)

var (
	frame uint64 = 0

	// debugging flags
	uiGrid = flag.Bool("ui-grid", false, "display UI grid lines")
)

func drawBackground() {
	// #4c6b73
	bottom := rl.NewColor(76, 107, 115, 255)
	// #313862
	via := rl.NewColor(49, 56, 98, 255)
	// #011b32
	top := rl.NewColor(1, 27, 50, 255)
	// 288.0
	viaPosition := float32(platform.WindowHeight) * 0.6
	// viaPosition := float32(platform.WindowHeight) * float32(frame%600) / 100.0
	rl.DrawRectangleGradientEx(
		rl.NewRectangle(0, 0, platform.WindowWidth, viaPosition),
		top,
		via,
		via,
		top,
	)
	rl.DrawRectangleGradientEx(
		rl.NewRectangle(0, viaPosition, platform.WindowWidth, platform.WindowHeight-viaPosition),
		via,
		bottom,
		bottom,
		via,
	)
}

func drawLayoutGrid() {
	color := rl.Red

	// outer margins
	margin := int32(20)
	// left
	rl.DrawLine(0, margin, platform.WindowWidth, margin, color)
	// right
	rl.DrawLine(0, platform.WindowHeight-margin, platform.WindowWidth, platform.WindowHeight-margin, color)
	// top
	rl.DrawLine(margin, 0, margin, platform.WindowHeight, color)
	// bottom
	rl.DrawLine(platform.WindowWidth-margin, 0, platform.WindowWidth-margin, platform.WindowHeight, color)

	// clock / date
	clockWidth := int32(800)
	rl.DrawLine(clockWidth-margin, 0, clockWidth-margin, platform.WindowHeight, color)
	rl.DrawLine(clockWidth, 0, clockWidth, platform.WindowHeight, color)
	rl.DrawLine(clockWidth+margin, 0, clockWidth+margin, platform.WindowHeight, color)
	rl.DrawLine(400, 0, 400, platform.WindowHeight, color)
	rl.DrawLine(0, 100, platform.WindowWidth, 100, color)
	rl.DrawLine(0, platform.WindowHeight-100, platform.WindowWidth, platform.WindowHeight-100, color)

	// right carousel
	carouselWidth := platform.WindowWidth - clockWidth - (2 * margin)
	rl.DrawLine(clockWidth+margin+(carouselWidth/2), 0, clockWidth+margin+(carouselWidth/2), platform.WindowHeight, color)
	rl.DrawLine(clockWidth+margin+(carouselWidth/4), 0, clockWidth+margin+(carouselWidth/4), platform.WindowHeight, color)
	rl.DrawLine(clockWidth+margin+(carouselWidth*3/4), 0, clockWidth+margin+(carouselWidth*3/4), platform.WindowHeight, color)
}

func drawTemperature(font rl.Font) {
	pos := rl.NewVector2(800, 0)
	rl.DrawTextEx(font, "88°", pos, float32(font.BaseSize), -16.0, rl.White)
	// TODO draw images
	// TODO test animated images
}

func RunForever() {
	rl.InitWindow(platform.ScreenWidth, platform.ScreenHeight, "clock")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	texture := rl.LoadRenderTexture(platform.WindowWidth, platform.WindowHeight)
	fontDefault := rl.LoadFontEx("assets/fonts/Oswald-Regular.ttf", 500, []rune{'8', '8', '°'})

	clock := widgets.NewClock()
	clock.Initialize()

	for !rl.WindowShouldClose() {
		frame += 1
		ctx := context.WithValue(context.Background(), widgets.KeyFrame, frame)
		rl.BeginDrawing()

		rl.BeginTextureMode(texture)
		drawBackground()
		if *uiGrid {
			drawLayoutGrid()
		}
		clock.Render(ctx, 0, 0, platform.ClockWidth, platform.WindowHeight)
		drawTemperature(fontDefault)
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

		rl.DrawTexturePro(texture.Texture, src, dst, rl.NewVector2(0, 0), rotation, rl.White)

		rl.EndDrawing()
	}
}
