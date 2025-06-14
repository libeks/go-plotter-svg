package fonts

import (
	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
)

type option struct {
	fontFile string
}

type textOption func(option) option

func withFontOption(path string) textOption {
	return func(o option) option {
		o.fontFile = path
		return o
	}
}

func RenderText(b box.Box, text string, textOptions ...textOption) []lines.LineLike {
	o := option{
		fontFile: "C:/Windows/Fonts/bahnschrift.ttf",
	}
	for _, opt := range textOptions {
		o = opt(o)
	}
	f, err := LoadFont(o.fontFile)
	if err != nil {
		panic(err)
	}
	curves := []lines.LineLike{}
	for _, ch := range text {
		// TODO: Properly space characters, they're rendered on top of each other at this point
		glyph, err := f.LoadGlyph(ch)
		if err != nil {
			panic(err)
		}
		c := glyph.GetCurves(b)
		curves = append(curves, c...)
	}
	return curves
}
