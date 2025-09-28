package widgets

import (
	"image/color"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type animatedIcon struct {
	img          rl.Image
	texture      rl.Texture2D
	frameCurrent int32
	frameTotal   int32
}

func NewAnimatedIcon(filename string) animatedIcon {
	a := animatedIcon{}
	a.LoadImage(filename)
	return a
}

func (a *animatedIcon) LoadImage(filename string) {
	if a.img.Data != nil {
		rl.UnloadImage(&a.img)
		a.frameTotal = 0
		a.frameCurrent = 0
	}
	a.img = *rl.LoadImageAnim(filename, &a.frameTotal)
	a.texture = rl.LoadTextureFromImage(&a.img)
}

func (a *animatedIcon) RenderFrame() {
	if a.frameTotal <= 1 {
		return
	}
	a.frameCurrent += 1
	if a.frameCurrent >= a.frameTotal {
		a.frameCurrent = 0
	}
	dataOffset := a.img.Width * a.img.Height * 4 * a.frameCurrent
	imgSize := a.img.Width * a.img.Height
	rl.UpdateTexture(
		a.texture,
		unsafe.Slice((*color.RGBA)(unsafe.Pointer(uintptr(a.img.Data)+uintptr(dataOffset))), imgSize),
	)
}

func (a *animatedIcon) Texture() rl.Texture2D {
	return a.texture
}
