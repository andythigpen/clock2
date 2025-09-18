package widgets

import (
	"context"

	"github.com/andythigpen/clock2/pkg/platform"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type carouselState int

const (
	carouselStateFadeIn carouselState = iota
	carouselStateNormal
	carouselStateFadeOut
)

type carousel struct {
	baseWidget
	widgets         []Widget
	index           int
	state           carouselState
	transitionStart uint64
	transitionEnd   uint64
}

func (c *carousel) currentWidget() Widget {
	return c.widgets[c.index]
}

func (c *carousel) RenderTexture(ctx context.Context) {
	frame := ctx.Value(KeyFrame).(uint64)

	if frame >= c.transitionEnd {
		c.transitionStart = frame
		switch c.state {
		case carouselStateFadeIn:
			c.state = carouselStateNormal
			c.transitionEnd = frame + (15 * platform.FPS)
		case carouselStateNormal:
			c.state = carouselStateFadeOut
			c.transitionEnd = frame + platform.FPS
		case carouselStateFadeOut:
			c.state = carouselStateFadeIn
			c.transitionEnd = frame + platform.FPS
			c.index += 1
			if c.index >= len(c.widgets) {
				c.index = 0
			}
		}
	}

	widget := c.currentWidget()
	widget.RenderTexture(ctx)

	rl.BeginTextureMode(c.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)
	texture := widget.Texture()
	var alpha uint8
	switch c.state {
	case carouselStateNormal:
		alpha = 255
	case carouselStateFadeIn:
		alpha = uint8(rl.Remap(float32(frame), float32(c.transitionStart), float32(c.transitionEnd), 0, 255))
	case carouselStateFadeOut:
		alpha = uint8(rl.Remap(float32(frame), float32(c.transitionStart), float32(c.transitionEnd), 255, 0))
	}
	rl.DrawTexturePro(texture,
		rl.NewRectangle(0, 0, float32(texture.Width), -float32(texture.Height)),
		rl.NewRectangle(0, 0, float32(texture.Width), float32(texture.Height)),
		rl.NewVector2(0, 0),
		0.0,
		rl.NewColor(255, 255, 255, alpha),
	)
}

func (c *carousel) ShouldDisplay() bool {
	return true
}

func NewCarousel(x, y float32, width, height int32, widgets ...Widget) Widget {
	if len(widgets) == 0 {
		panic("at least one widget is required")
	}
	return &carousel{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
			x:       x,
			y:       y,
		},
		widgets: widgets,
		index:   0,
	}
}
