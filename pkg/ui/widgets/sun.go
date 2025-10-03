package widgets

import (
	"context"
	"time"

	"github.com/andythigpen/clock2/pkg/services"
	"github.com/andythigpen/clock2/pkg/ui/widgets/fonts"
	"github.com/andythigpen/clock2/pkg/ui/widgets/icons"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type sun struct {
	baseWidget
	svc         *services.HomeAssistantService
	rising      bool
	fontClock   rl.Font
	iconRising  icons.AnimatedIcon
	iconSetting icons.AnimatedIcon
	nextRising  time.Time
	nextSetting time.Time
}

var _ Fetcher = (*sun)(nil)

func (s *sun) FetchData(ctx context.Context) {
	sun := s.svc.GetSun()
	s.rising = sun.Attributes.Rising
	s.nextRising = sun.Attributes.NextRising.Local()
	s.nextSetting = sun.Attributes.NextSetting.Local()
}

func (s *sun) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(s.texture)
	defer rl.EndTextureMode()
	rl.ClearBackground(rl.Blank)

	var textTime string
	if s.rising {
		x := s.texture.Texture.Width/4 - (s.iconRising.Width() / 2)
		y := s.texture.Texture.Height/2 - (s.iconRising.Height() / 2)
		s.iconRising.RenderFrame(float32(x), float32(y))
		textTime = s.nextRising.Format("03:04")
	} else {
		x := s.texture.Texture.Width/4 - (s.iconSetting.Width() / 2)
		y := s.texture.Texture.Height/2 - (s.iconSetting.Height() / 2)
		s.iconSetting.RenderFrame(float32(x), float32(y))
		textTime = s.nextSetting.Format("03:04")
	}

	spacing := float32(-8.0)
	textSize := rl.MeasureTextEx(s.fontClock, textTime, float32(s.fontClock.BaseSize), spacing)
	textX := float32(s.texture.Texture.Width)/2 + spacing
	textY := float32(s.texture.Texture.Height)/2 - textSize.Y/2
	rl.DrawTextEx(
		s.fontClock,
		textTime,
		rl.NewVector2(textX, textY),
		float32(s.fontClock.BaseSize),
		spacing,
		rl.White,
	)
}

func (s *sun) ShouldDisplay() bool {
	twoHoursFromNow := time.Now().Add(2 * time.Hour)
	return s.nextRising.Before(twoHoursFromNow) || s.nextSetting.Before(twoHoursFromNow)
}

func (s *sun) LoadAssets() {
	if s.rising {
		s.iconRising.LoadAssets()
	} else {
		s.iconSetting.LoadAssets()
	}
}

func (s *sun) UnloadAssets() {
	if s.rising {
		s.iconRising.UnloadAssets()
	} else {
		s.iconSetting.UnloadAssets()
	}
}

func NewSun(width, height int32, svc *services.HomeAssistantService) Widget {
	return &sun{
		baseWidget:  newBaseWidget(0, 0, width, height),
		svc:         svc,
		iconRising:  icons.NewAnimatedIcon(icons.IconSunrise),
		iconSetting: icons.NewAnimatedIcon(icons.IconSunset),
		fontClock:   fonts.Cache.Load(fonts.FontBebasNeue, 340),
	}
}
