package widgets

import (
	"context"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Widget interface {
	// Call drawing functions to draw the Widget onto a texture
	RenderTexture(ctx context.Context)
	// Return the texture that will be drawn to the screen
	Texture() rl.Texture2D
	// Return true to display the widget, false otherwise
	ShouldDisplay() bool
	GetX() float32
	GetY() float32
}

type Fetcher interface {
	// Called prior to the rendering loop so that widgets can fetch information for display
	FetchData(ctx context.Context)
}

type baseWidget struct {
	texture rl.RenderTexture2D
	x       float32
	y       float32
}

func (b *baseWidget) Texture() rl.Texture2D {
	return b.texture.Texture
}

func (b *baseWidget) GetX() float32 {
	return b.x
}

func (b *baseWidget) GetY() float32 {
	return b.y
}
