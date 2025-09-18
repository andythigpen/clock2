package widgets

import (
	"context"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type background struct {
	baseWidget
}

func (b *background) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(b.texture)
	defer rl.EndTextureMode()

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
