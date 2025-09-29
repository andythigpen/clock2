package widgets

import (
	"context"
	"flag"

	"github.com/andythigpen/clock2/pkg/platform"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	uiTestCarouselWidget = flag.Int("ui-test-carousel-widget", -1, "index of a carousel widget to test")
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
	shouldAdvance   bool
}

var _ Fetcher = (*carousel)(nil)

func (c *carousel) currentWidget() Widget {
	return c.widgets[c.index]
}

func (c *carousel) advanceTransition(ctx context.Context, frame uint64) {
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
		if c.shouldAdvance {
			for range c.widgets {
				c.index += 1
				if c.index >= len(c.widgets) {
					c.index = 0
				}
				if f, ok := (c.widgets[c.index]).(Fetcher); ok {
					f.FetchData(ctx)
				}
				if c.widgets[c.index].ShouldDisplay() {
					break
				}
			}
		}
	}
}

func (c *carousel) FetchData(ctx context.Context) {
	widget := c.currentWidget()
	if f, ok := widget.(Fetcher); ok {
		f.FetchData(ctx)
	}

	frame := ctx.Value(KeyFrame).(uint64)
	if frame >= c.transitionEnd {
		c.advanceTransition(ctx, frame)
	}
}

func (c *carousel) RenderTexture(ctx context.Context) {
	frame := ctx.Value(KeyFrame).(uint64)

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
	c := &carousel{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
			x:       x,
			y:       y,
		},
		widgets: widgets,
		index:   0,
	}
	if *uiTestCarouselWidget >= 0 && *uiTestCarouselWidget < len(c.widgets) {
		c.index = *uiTestCarouselWidget
		c.shouldAdvance = false
	} else {
		c.shouldAdvance = true
	}
	return c
}
