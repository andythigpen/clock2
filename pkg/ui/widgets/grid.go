package widgets

import (
	"context"
	"flag"

	"github.com/andythigpen/clock2/pkg/platform"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// debugging flags
	uiGrid = flag.Bool("ui-grid", false, "display UI grid lines")
)

type grid struct {
	baseWidget
}

func (g *grid) RenderTexture(ctx context.Context) {
	rl.BeginTextureMode(g.texture)
	defer rl.EndTextureMode()

	color := rl.Red
	colorThirds := rl.Green
	colorQuarters := rl.Purple
	width := g.texture.Texture.Width
	height := g.texture.Texture.Height

	// outer margins
	margin := int32(platform.Margin)
	// left
	rl.DrawLine(0, margin, width, margin, color)
	// right
	rl.DrawLine(0, height-margin, width, height-margin, color)
	// top
	rl.DrawLine(margin, 0, margin, height, color)
	// bottom
	rl.DrawLine(width-margin, 0, width-margin, height, color)

	// clock / date
	clockWidth := int32(800)
	rl.DrawLine(clockWidth-margin, 0, clockWidth-margin, height, color)
	rl.DrawLine(clockWidth, 0, clockWidth, height, color)
	rl.DrawLine(clockWidth+margin, 0, clockWidth+margin, height, color)
	rl.DrawLine(400, 0, 400, height, color)
	rl.DrawLine(0, 100, width, 100, color)
	rl.DrawLine(0, height-100, width, height-100, color)

	// right carousel
	carouselWidth := width - clockWidth - (2 * margin)
	rl.DrawLine(clockWidth+margin+(carouselWidth/2), 0, clockWidth+margin+(carouselWidth/2), height, color)
	rl.DrawLine(clockWidth+margin+(carouselWidth/4), 0, clockWidth+margin+(carouselWidth/4), height, colorQuarters)
	rl.DrawLine(clockWidth+margin+(carouselWidth*3/4), 0, clockWidth+margin+(carouselWidth*3/4), height, colorQuarters)
	rl.DrawLine(clockWidth+margin+(carouselWidth/3), 0, clockWidth+margin+(carouselWidth/3), height, colorThirds)
	rl.DrawLine(clockWidth+margin+(carouselWidth*2/3), 0, clockWidth+margin+(carouselWidth*2/3), height, colorThirds)
}

func (g *grid) ShouldDisplay() bool {
	return *uiGrid
}

func NewGrid(width, height int32) Widget {
	return &grid{
		baseWidget: newBaseWidget(0, 0, width, height),
	}
}
