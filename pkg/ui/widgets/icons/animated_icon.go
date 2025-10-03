package icons

import (
	"fmt"
	"log/slog"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/andythigpen/clock2/pkg/platform"
)

type iconData struct {
	cols int32
	rows int32
}

var (
	icons = map[IconType]iconData{
		IconClearDay:               {5, 5},
		IconClearNight:             {7, 8},
		IconCloudyDay:              {5, 8},
		IconCloudyNight:            {5, 8},
		IconCodeOrange:             {8, 8},
		IconCodeRed:                {8, 8},
		IconFogDay:                 {5, 7},
		IconFogNight:               {6, 7},
		IconHail:                   {6, 6},
		IconHumidity:               {13, 8},
		IconPartlyCloudyDay:        {5, 8},
		IconPartlyCloudyNight:      {6, 8},
		IconPressureHigh:           {28, 5},
		IconPressureLow:            {28, 5},
		IconRain:                   {6, 6},
		IconSleet:                  {6, 6},
		IconSnow:                   {6, 6},
		IconSunrise:                {5, 9},
		IconSunset:                 {5, 7},
		IconThunderstormsDay:       {5, 5},
		IconThunderstormsNight:     {6, 6},
		IconThunderstormsRainDay:   {5, 5},
		IconThunderstormsRainNight: {6, 5},
		IconWindy:                  {6, 9},
	}
)

type AnimatedIcon struct {
	iconType       IconType
	textures       []rl.Texture2D
	textureCurrent int32
	frameCounter   int32 // counts to the frame speed and resets
	frameCurrent   int32 // current frame
	frameSpeed     int32 // number of frames shown per second
	frameTotal     int32 // total number of frames in the sprite
	frameRec       rl.Rectangle
	cols           int32 // number of columns in the sprite sheet
	rows           int32 // number of rows in the sprite sheet
}

func NewAnimatedIcon(iconType IconType) AnimatedIcon {
	a := AnimatedIcon{frameSpeed: 20, frameTotal: 120}
	a.SetIconType(iconType)
	return a
}

func (a *AnimatedIcon) RenderFrame(x, y float32) {
	if a.frameTotal <= 1 {
		return
	}
	a.frameCounter += 1
	if a.frameCounter >= (platform.FPS / a.frameSpeed) {
		a.frameCounter = 0
		a.frameCurrent += 1
		a.textureCurrent = a.frameCurrent / (a.cols * a.rows)
		if a.frameCurrent >= a.frameTotal {
			a.frameCurrent = 0
			a.textureCurrent = 0
		}
		if int(a.textureCurrent) >= len(a.textures) {
			slog.Warn(
				"missing textures",
				"frameCurrent", a.frameCurrent,
				"frameTotal", a.frameTotal,
				"textureCurrent", a.textureCurrent,
				"total", len(a.textures),
				"rows", a.rows,
				"cols", a.cols,
				"iconType", a.iconType,
			)
			a.textureCurrent = 0
		}
		a.frameRec.X = float32(a.frameCurrent%a.cols) * a.frameRec.Width
		a.frameRec.Y = float32(a.frameCurrent%(a.cols*a.rows)/a.cols) * a.frameRec.Height
	}
	texture := a.textures[a.textureCurrent]
	rl.DrawTextureRec(texture, a.frameRec, rl.NewVector2(x, y), rl.White)
}

func (a *AnimatedIcon) Width() int32 {
	return int32(a.frameRec.Width)
}

func (a *AnimatedIcon) Height() int32 {
	return int32(a.frameRec.Height)
}

func (a *AnimatedIcon) SetIconType(iconType IconType) {
	if _, ok := icons[iconType]; ok {
		a.iconType = iconType
	}
}

func (a *AnimatedIcon) LoadAssets() {
	a.frameCurrent = 0
	a.textureCurrent = 0
	a.textures = make([]rl.Texture2D, 0)
	matches, _ := filepath.Glob(fmt.Sprintf("assets/icons/weather/sprites/%s-[0-9].png", a.iconType))
	slog.Info("LoadAssets before", "iconType", a.iconType, "matches", matches)
	for _, filename := range matches {
		slog.Info("LoadAssets", "filename", filename)
		a.textures = append(a.textures, rl.LoadTexture(filename))
	}
	d := icons[a.iconType]
	a.cols = d.cols
	a.rows = d.rows
	if a.cols > 0 && a.rows > 0 && len(a.textures) > 0 {
		texture := a.textures[0]
		a.frameRec = rl.NewRectangle(0, 0, float32(texture.Width/a.cols), float32(texture.Height/a.rows))
	}
}

func (a *AnimatedIcon) UnloadAssets() {
	for _, texture := range a.textures {
		rl.UnloadTexture(texture)
	}
	a.textures = nil
}
