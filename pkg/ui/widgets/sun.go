package widgets

import (
	"context"
	"time"

	"github.com/andythigpen/clock2/pkg/services"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type sun struct {
	baseWidget
	svc         *services.HomeAssistantService
	rising      bool
	fontClock   rl.Font
	iconRising  animatedIcon
	iconSetting animatedIcon
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

	var (
		texture  rl.Texture2D
		textTime string
	)
	if s.rising {
		s.iconRising.RenderFrame()
		texture = s.iconRising.Texture()
		textTime = s.nextRising.Format("03:04")
	} else {
		s.iconSetting.RenderFrame()
		texture = s.iconSetting.Texture()
		textTime = s.nextSetting.Format("03:04")
	}
	rl.DrawTexture(texture, 50, 0, rl.White)

	spacing := float32(-8.0)
	textSize := rl.MeasureTextEx(s.fontClock, textTime, float32(s.fontClock.BaseSize), spacing)
	textX := float32(s.texture.Texture.Width) / 2
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

func NewSun(width, height int32, svc *services.HomeAssistantService) Widget {
	return &sun{
		baseWidget: baseWidget{
			texture: rl.LoadRenderTexture(width, height),
		},
		svc:         svc,
		iconRising:  NewAnimatedIcon(getAssetIconPath("sunrise")),
		iconSetting: NewAnimatedIcon(getAssetIconPath("sunset")),
		fontClock:   rl.LoadFontEx("assets/fonts/BebasNeue-Regular.ttf", 340, nil),
	}
}
