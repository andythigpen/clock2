package widgets

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/platform"
)

type animatedIcon struct {
	filename     string
	texture      rl.Texture2D
	frameCounter int32 // counts to the frame speed and resets
	frameCurrent int32 // current frame
	frameSpeed   int32 // number of frames shown per second
	frameTotal   int32 // total number of frames in the sprite
	frameRec     rl.Rectangle
	cols         int32 // number of columns in the sprite sheet
	rows         int32 // number of rows in the sprite sheet
}

type animatedIconOption func(*animatedIcon)

func WithCols(cols int32) animatedIconOption {
	return func(a *animatedIcon) {
		a.cols = cols
	}
}

func WithRows(rows int32) animatedIconOption {
	return func(a *animatedIcon) {
		a.rows = rows
	}
}

func WithTotalFrames(total int32) animatedIconOption {
	return func(a *animatedIcon) {
		a.frameTotal = total
	}
}

func NewAnimatedIcon(filename string, opts ...animatedIconOption) animatedIcon {
	a := animatedIcon{filename: filename, frameSpeed: 20, frameTotal: 120, cols: 30, rows: 4}
	for _, o := range opts {
		o(&a)
	}
	return a
}

func (a *animatedIcon) RenderFrame(x, y float32) {
	if a.frameTotal <= 1 {
		return
	}
	a.frameCounter += 1
	if a.frameCounter >= (platform.FPS / a.frameSpeed) {
		a.frameCounter = 0
		a.frameCurrent += 1
		if a.frameCurrent >= a.frameTotal {
			a.frameCurrent = 0
		}
		a.frameRec.X = float32(a.frameCurrent%a.cols) * a.frameRec.Width
		a.frameRec.Y = float32(a.frameCurrent%a.frameTotal/a.cols) * a.frameRec.Height
	}

	rl.DrawTextureRec(a.texture, a.frameRec, rl.NewVector2(x, y), rl.White)
}

func (a *animatedIcon) Width() int32 {
	return int32(a.frameRec.Width)
}

func (a *animatedIcon) Height() int32 {
	return int32(a.frameRec.Height)
}

func (a *animatedIcon) SetFilename(filename string) {
	a.filename = filename
}

func (a *animatedIcon) LoadAssets() {
	a.frameCurrent = 0
	a.texture = rl.LoadTexture(a.filename)
	if a.cols > 0 && a.rows > 0 {
		a.frameRec = rl.NewRectangle(0, 0, float32(a.texture.Width/a.cols), float32(a.texture.Height/a.rows))
	}
}

func (a *animatedIcon) UnloadAssets() {
	rl.UnloadTexture(a.texture)
}
