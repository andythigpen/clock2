package widgets

import (
	"context"
	"flag"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	uiTestBackground = flag.Bool("ui-test-background", false, "when true, cycle through background states")
)

type background struct {
	baseWidget
}

func (b *background) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(b.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.White)

	if *uiTestBackground {
		frame := ctx.Value(KeyFrame).(uint64)
		mins := frame % 1800
		var (
			bottom, via, top rl.Color
			viaPosition      float32
		)
		if mins < 360 {
			// morning
			topStart := rl.Black
			topEnd := rl.NewColor(1, 27, 50, 255)
			top = rl.ColorLerp(topStart, topEnd, rl.Remap(float32(mins), 0, 360, 0.0, 1.0))
			bottom = rl.NewColor(21, 13, 13, 255)
			viaStart := topStart
			viaEnd := rl.NewColor(49, 56, 98, 255)
			via = rl.ColorLerp(viaStart, viaEnd, rl.Remap(float32(mins), 0, 360, 0.0, 1.0))
			viaPosition = float32(b.texture.Texture.Height)
		} else if mins >= 360 && mins < 720 {
			// sunrise
			top = rl.NewColor(1, 27, 50, 255)
			bottomStart := rl.NewColor(115, 76, 103, 255)
			bottomEnd := rl.NewColor(76, 107, 115, 255)
			bottom = rl.ColorLerp(bottomStart, bottomEnd, rl.Remap(float32(mins), 360, 720, 0.0, 1.0))
			via = rl.NewColor(49, 56, 98, 255)
			viaPosition = rl.Remap(float32(mins), 360, 720, float32(b.texture.Texture.Height), float32(b.texture.Texture.Height)*0.5)
		} else if mins >= 720 && mins < 1080 {
			// day
			top = rl.NewColor(1, 27, 50, 255)
			bottom = rl.NewColor(76, 107, 115, 255)
			via = rl.NewColor(49, 56, 98, 255)
			viaPosition = float32(b.texture.Texture.Height) * 0.5
		} else if mins >= 1080 && mins < 1440 {
			// sunset
			top = rl.NewColor(1, 27, 50, 255)
			bottomEnd := rl.NewColor(115, 76, 103, 255)
			bottomStart := rl.NewColor(76, 107, 115, 255)
			bottom = rl.ColorLerp(bottomStart, bottomEnd, rl.Remap(float32(mins), 1080, 1440, 0.0, 1.0))
			via = rl.NewColor(49, 56, 98, 255)
			viaPosition = rl.Remap(float32(mins), 1080, 1440, float32(b.texture.Texture.Height)*0.5, float32(b.texture.Texture.Height))
		} else {
			// night
			topEnd := rl.Black
			topStart := rl.NewColor(1, 27, 50, 255)
			top = rl.ColorLerp(topStart, topEnd, rl.Remap(float32(mins), 1440, 1800, 0.0, 1.0))
			bottom = rl.NewColor(21, 13, 13, 255)
			viaEnd := topEnd
			viaStart := rl.NewColor(49, 56, 98, 255)
			via = rl.ColorLerp(viaStart, viaEnd, rl.Remap(float32(mins), 1440, 1800, 0.0, 1.0))
			viaPosition = float32(b.texture.Texture.Height)
		}
		rl.DrawRectangleGradientEx(
			rl.NewRectangle(0, 0, float32(b.texture.Texture.Width), viaPosition),
			top,
			via,
			via,
			top,
		)
		rl.DrawRectangleGradientEx(
			rl.NewRectangle(
				0,
				viaPosition,
				float32(b.texture.Texture.Width),
				float32(b.texture.Texture.Height)-viaPosition,
			),
			via,
			bottom,
			bottom,
			via,
		)

		return
	}

	// #4c6b73
	bottom := rl.NewColor(76, 107, 115, 255)
	// #313862
	via := rl.NewColor(49, 56, 98, 255)
	// #011b32
	top := rl.NewColor(1, 27, 50, 255)
	// 288.0
	viaPosition := float32(b.texture.Texture.Height) * 0.6
	// viaPosition := float32(platform.WindowHeight) * float32(frame%600) / 100.0
	rl.DrawRectangleGradientEx(
		rl.NewRectangle(0, 0, float32(b.texture.Texture.Width), viaPosition),
		top,
		via,
		via,
		top,
	)
	rl.DrawRectangleGradientEx(
		rl.NewRectangle(
			0,
			viaPosition,
			float32(b.texture.Texture.Width),
			float32(b.texture.Texture.Height)-viaPosition,
		),
		via,
		bottom,
		bottom,
		via,
	)
}

func (b *background) ShouldDisplay() bool {
	return true
}

func NewBackground(x, y float32, width, height int32) Widget {
	return &background{
		baseWidget{
			texture: rl.LoadRenderTexture(width, height),
			x:       x,
			y:       y,
		},
	}
}
