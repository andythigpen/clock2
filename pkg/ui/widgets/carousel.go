package widgets

import (
	"context"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type carousel struct {
	x       float32
	y       float32
	widgets []Widget
	index   int
}

func (c *carousel) currentWidget() Widget {
	return c.widgets[c.index]
}

func (c *carousel) RenderTexture(ctx context.Context) {
	widget := c.currentWidget()
	widget.RenderTexture(ctx)
}

func (c *carousel) ShouldDisplay() bool {
	return true
}

func (c *carousel) Texture() rl.Texture2D {
	widget := c.currentWidget()
	return widget.Texture()
}

func (c *carousel) GetX() float32 {
	return c.x
}

func (c *carousel) GetY() float32 {
	return c.y
}

func NewCarousel(x, y float32, widgets ...Widget) Widget {
	if len(widgets) == 0 {
		panic("at least one widget is required")
	}
	return &carousel{
		widgets: widgets,
		index:   0,
		x:       x,
		y:       y,
	}
}
