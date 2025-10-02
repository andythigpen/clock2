package widgets

import (
	"context"
	"flag"
	"time"

	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	uiTestBackground = flag.Bool("ui-test-background", false, "when true, cycle through background states")
)

func minsAfterMidnight(t time.Time) int {
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int(t.Sub(midnight).Minutes())
}

type background struct {
	baseWidget
	svc *services.HomeAssistantService
}

var (
	colorDayTop    = rl.NewColor(1, 27, 50, 255)
	colorDayVia    = rl.NewColor(49, 56, 98, 255)
	colorDayBottom = rl.NewColor(76, 107, 115, 255)

	colorDawnDuskBottom      = rl.NewColor(21, 13, 13, 255)
	colorSunriseSunsetBottom = rl.NewColor(115, 76, 103, 255)
)

func (b *background) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(b.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.White)

	var (
		minCurrent int
	)

	now := time.Now().Local()
	sun := b.svc.GetSun()
	minTotal := 1440
	minSunrise := minsAfterMidnight(sun.Attributes.NextRising.Local())
	minAfterSunrise := minSunrise + 90
	minSunset := minsAfterMidnight(sun.Attributes.NextSetting.Local())
	minBeforeSunset := minSunset - 90

	if *uiTestBackground {
		frame := ctx.Value(KeyFrame).(uint64)
		minCurrent = int(frame % uint64(minTotal))
	} else {
		minCurrent = minsAfterMidnight(now)
	}

	var (
		bottom, via, top rl.Color
		viaPosition      float32
	)

	if minCurrent < minSunrise {
		// morning
		topStart := rl.Black
		topEnd := colorDayTop
		top = rl.ColorLerp(topStart, topEnd, rl.Remap(float32(minCurrent), 0, float32(minSunrise), 0.0, 1.0))
		bottom = colorDawnDuskBottom
		viaStart := topStart
		viaEnd := colorDayVia
		via = rl.ColorLerp(viaStart, viaEnd, rl.Remap(float32(minCurrent), 0, float32(minSunrise), 0.0, 1.0))
		viaPosition = float32(b.texture.Texture.Height)
	} else if minCurrent >= minSunrise && minCurrent < minAfterSunrise {
		// sunrise
		top = colorDayTop
		bottomStart := colorSunriseSunsetBottom
		bottomEnd := colorDayBottom
		bottom = rl.ColorLerp(bottomStart, bottomEnd, rl.Remap(float32(minCurrent), float32(minSunrise), float32(minAfterSunrise), 0.0, 1.0))
		via = colorDayVia
		viaPosition = rl.Remap(float32(minCurrent), float32(minSunrise), float32(minAfterSunrise), float32(b.texture.Texture.Height), float32(b.texture.Texture.Height)*0.5)
	} else if minCurrent >= minAfterSunrise && minCurrent < minBeforeSunset {
		// day
		top = colorDayTop
		bottom = colorDayBottom
		via = colorDayVia
		viaPosition = float32(b.texture.Texture.Height) * 0.5
	} else if minCurrent >= minBeforeSunset && minCurrent < minSunset {
		// sunset
		top = colorDayTop
		bottomEnd := colorSunriseSunsetBottom
		bottomStart := colorDayBottom
		bottom = rl.ColorLerp(bottomStart, bottomEnd, rl.Remap(float32(minCurrent), float32(minBeforeSunset), float32(minSunset), 0.0, 1.0))
		via = colorDayVia
		viaPosition = rl.Remap(float32(minCurrent), float32(minBeforeSunset), float32(minSunset), float32(b.texture.Texture.Height)*0.5, float32(b.texture.Texture.Height))
	} else {
		// night
		topEnd := rl.Black
		topStart := colorDayTop
		top = rl.ColorLerp(topStart, topEnd, rl.Remap(float32(minCurrent), float32(minSunset), float32(minTotal), 0.0, 1.0))
		bottom = colorDawnDuskBottom
		viaEnd := topEnd
		viaStart := colorDayVia
		via = rl.ColorLerp(viaStart, viaEnd, rl.Remap(float32(minCurrent), float32(minSunset), float32(minTotal), 0.0, 1.0))
		viaPosition = float32(b.texture.Texture.Height)
	}
	rl.DrawRectangleGradientEx(
		rl.NewRectangle(0, 0, float32(b.texture.Texture.Width), viaPosition),
		top,
		via,
		via,
		top,
	)
	rl.DrawRectangleGradientEx(
		rl.NewRectangle(
			0,
			viaPosition,
			float32(b.texture.Texture.Width),
			float32(b.texture.Texture.Height)-viaPosition,
		),
		via,
		bottom,
		bottom,
		via,
	)
}

func (b *background) ShouldDisplay() bool {
	return true
}

func NewBackground(x, y float32, width, height int32, svc *services.HomeAssistantService) Widget {
	return &background{
		baseWidget: newBaseWidget(x, y, width, height),
		svc:        svc,
	}
}
