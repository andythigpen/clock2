package widgets

import (
	"context"
)

type Widget interface {
	// Load any fonts, textures, etc. when called on startup
	Initialize()
	// Call drawing functions to draw the Widget within intialized width and height
	// TODO: set the frame as a context variable before calling
	Render(ctx context.Context, x, y, width, height float32)
	ShouldDisplay() bool
}
