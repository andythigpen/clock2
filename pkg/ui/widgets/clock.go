package widgets

import (
	"context"
	"flag"
	"fmt"
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
	if *uiTestClock {
		currentTime := ctx.Value(KeyFrame).(uint64) % 5184000 // 60 FPS * 60s * 60m * 24h
		hour := (currentTime / 60 % 12) + 1
		minute := currentTime / 20 % 60
		display = fmt.Sprintf("%02d:%02d", hour, minute)
	} else {
		display = time.Now().Format("03:04")
	}
	pos := rl.NewVector2(platform.Margin, 2*platform.Margin)
	rl.DrawTextEx(c.fontClock, display, pos, float32(c.fontClock.BaseSize), -8.0, rl.White)

	v := rl.MeasureTextEx(c.fontDate, "Mon Sep 15", float32(c.fontDate.BaseSize), 0.0)
	spacing := (float32(c.texture.Texture.Width) - platform.Margin - v.X) / 2.0
	pos = rl.NewVector2(platform.Margin+spacing, -40)
	rl.DrawTextEx(c.fontDate, "Mon Sep 15", pos, float32(c.fontDate.BaseSize), 0.0, rl.White)
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
