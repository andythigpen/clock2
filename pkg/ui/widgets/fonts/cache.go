package fonts

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type fontType string
type fontVariation string

const (
	FontOswald    fontType = "Oswald"
	FontBebasNeue fontType = "BebasNeue"
	FontMoulpali  fontType = "Moulpali"

	FontVariationRegular fontVariation = "Regular"
	FontVariationBold    fontVariation = "Bold"
)

type fontCache map[string]map[int32]rl.Font

type fontOptions struct {
	variation fontVariation
}

type fontOption func(*fontOptions)

func WithVariation(variation fontVariation) fontOption {
	return func(o *fontOptions) {
		o.variation = variation
	}
}

func (f *fontCache) Load(name fontType, size int32, options ...fontOption) rl.Font {
	var (
		m  map[int32]rl.Font
		ok bool
	)

	opts := fontOptions{variation: "Regular"}
	for _, o := range options {
		o(&opts)
	}

	key := fmt.Sprintf("%s-%s", name, opts.variation)
	if m, ok = (*f)[key]; !ok {
		(*f)[key] = make(map[int32]rl.Font)
		m = (*f)[key]
	}
	if font, ok := m[size]; ok {
		return font
	}
	path := GetAssetFontPath(name, options...)
	m[size] = rl.LoadFontEx(path, size, nil)
	return m[size]
}

func (f *fontCache) Unload() {
	for _, m := range *f {
		for _, font := range m {
			rl.UnloadFont(font)
		}
	}
	*f = make(fontCache)
}

var Cache = fontCache{}

func GetAssetFontPath(name fontType, options ...fontOption) string {
	opts := fontOptions{variation: "Regular"}
	for _, o := range options {
		o(&opts)
	}
	return path.Join("assets/fonts", fmt.Sprintf("%s-%s.ttf", name, opts.variation))
}
