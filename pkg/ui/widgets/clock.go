package widgets

import (
	"context"
	"flag"
	"time"

	"github.com/andythigpen/clock2/pkg/platform"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var uiTestClock = flag.Bool("ui-test-clock", false, "when true, cycle through all numbers quickly")

type clock struct {
	baseWidget
	fontClock rl.Font
	fontDate  rl.Font
}

func (c *clock) ShouldDisplay() bool {
	return true // always displayed
}

func (c *clock) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(c.texture)
	defer rl.EndTextureMode()

	rl.ClearBackground(rl.Blank)

	var display string
	now := time.Now()
	if *uiTestClock {
		frame := ctx.Value(KeyFrame).(uint64)
		now = now.Add(time.Duration(frame) * time.Second)
	}
	display = now.Format("03:04")
	pos := rl.NewVector2(platform.Margin, 2*platform.Margin)
	rl.DrawTextEx(c.fontClock, display, pos, float32(c.fontClock.BaseSize), -8.0, rl.White)

	display = now.Format("Mon Jan _2")
	v := rl.MeasureTextEx(c.fontDate, display, float32(c.fontDate.BaseSize), 0.0)
	spacing := (float32(c.texture.Texture.Width) - platform.Margin - v.X) / 2.0
	pos = rl.NewVector2(platform.Margin+spacing, -40)
	rl.DrawTextEx(c.fontDate, display, pos, float32(c.fontDate.BaseSize), 0.0, rl.White)
}

func NewClock(x, y float32, width, height int32) Widget {
	runes := []rune{}
	for r := '0'; r <= ':'; r++ {
		runes = append(runes, r)
	}
	return &clock{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
			x:       x,
			y:       y,
		},
		fontClock: rl.LoadFontEx("assets/fonts/BebasNeue-Regular.ttf", 540, runes),
		fontDate:  rl.LoadFontEx("assets/fonts/Moulpali-Regular.ttf", 220, nil),
	}
}
