package widgets

import (
	"context"
	"flag"

	"github.com/andythigpen/clock2/pkg/platform"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	uiTestCarouselWidget  = flag.Int("ui-test-carousel-widget", -1, "index of a carousel widget to test")
	uiTestCarouselSeconds = flag.Int("ui-test-carousel-secs", 15, "number of seconds to display a widget before advancing")
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
	indexNext       int
	state           carouselState
	transitionFrame int
	transitionEnd   int
	shouldAdvance   bool
}

var _ Fetcher = (*carousel)(nil)

func (c *carousel) currentWidget() Widget {
	return c.widgets[c.index]
}

func (c *carousel) getNextIndex() int {
	if !c.shouldAdvance {
		return c.index
	}

	idx := c.index
	for range c.widgets {
		idx += 1
		if idx >= len(c.widgets) {
			idx = 0
		}
		if c.widgets[idx].ShouldDisplay() {
			break
		}
	}
	return idx
}

func (c *carousel) advanceTransition() {
	c.transitionFrame = 0
	switch c.state {
	case carouselStateFadeIn:
		c.state = carouselStateNormal
		c.transitionEnd = (*uiTestCarouselSeconds) * platform.FPS
	case carouselStateNormal:
		c.state = carouselStateFadeOut
		c.transitionEnd = 1 * platform.FPS
	case carouselStateFadeOut:
		c.state = carouselStateFadeIn
		c.transitionEnd = 1 * platform.FPS
		widget := c.currentWidget()
		if widget, ok := widget.(Loader); ok {
			widget.UnloadAssets()
		}
		c.indexNext = c.getNextIndex()
		c.index = c.indexNext
		widget = c.currentWidget()
		if widget, ok := widget.(Loader); ok {
			widget.LoadAssets()
		}
	}
}

func (c *carousel) FetchData(ctx context.Context) {
	for idx, widget := range c.widgets {
		if idx == c.index && *uiTestCarouselWidget == -1 {
			// skip the current widget so there are no abrupt changes if the data changes
			continue
		}
		if widget, ok := widget.(Fetcher); ok {
			widget.FetchData(ctx)
		}
	}
}

func (c *carousel) RenderTexture(ctx context.Context) {
	c.transitionFrame += 1
	if c.transitionFrame >= c.transitionEnd {
		c.advanceTransition()
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
		alpha = uint8(rl.Remap(float32(c.transitionFrame), 0, float32(c.transitionEnd), 0, 255))
	case carouselStateFadeOut:
		alpha = uint8(rl.Remap(float32(c.transitionFrame), 0, float32(c.transitionEnd), 255, 0))
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

func (c *carousel) LoadAssets() {
	widget := c.currentWidget()
	if widget, ok := widget.(Loader); ok {
		widget.LoadAssets()
	}
}

func (c *carousel) UnloadAssets() {
	widget := c.currentWidget()
	if widget, ok := widget.(Loader); ok {
		widget.UnloadAssets()
	}
	if c.index != c.indexNext {
		if widget, ok := c.widgets[c.indexNext].(Loader); ok {
			widget.UnloadAssets()
		}
	}
}

func (c *carousel) Unload() {
	c.baseWidget.Unload()
	for _, widget := range c.widgets {
		widget.Unload()
	}
}

func NewCarousel(x, y float32, width, height int32, widgets ...Widget) Widget {
	if len(widgets) == 0 {
		panic("at least one widget is required")
	}
	c := &carousel{
		baseWidget: newBaseWidget(x, y, width, height),
		widgets:    widgets,
		index:      0,
	}
	if *uiTestCarouselWidget >= 0 && *uiTestCarouselWidget < len(c.widgets) {
		c.index = *uiTestCarouselWidget
		c.shouldAdvance = false
	} else {
		c.shouldAdvance = true
	}
	return c
}
