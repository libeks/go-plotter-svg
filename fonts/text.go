package fonts

import (
	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
)

type option struct {
	fontFile string
}

// TextRenderer contains everything that you'd need to render text to a Scene
type TextRender struct {
	Text        string
	BoundingBox box.Box
	CharBoxes   []box.Box
	CharCurves  []lines.LineLike
	CharPoints  []lines.LineLike
}

type textOption func(option) option

func withFontOption(path string) textOption {
	return func(o option) option {
		o.fontFile = path
		return o
	}
}

func RenderText(b box.Box, text string, textOptions ...textOption) TextRender {
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
	controlPoints := []lines.LineLike{}
	curves := []lines.LineLike{}
	for _, ch := range text {
		// TODO: Properly space characters, they're rendered on top of each other at this point
		glyph, err := f.LoadGlyph(ch)
		if err != nil {
			panic(err)
		}
		c := glyph.GetCurves(b)
		curves = append(curves, c...)
		for _, pt := range glyph.GetControlPoints(b) {
			if pt.OnLine {
				controlPoints = append(controlPoints, objects.Circle{Center: pt.Point, Radius: 30})
			} else {
				controlPoints = append(controlPoints, objects.Circle{Center: pt.Point, Radius: 15})
			}
		}
	}

	return TextRender{
		Text:       text,
		CharCurves: curves,
		CharPoints: controlPoints,
	}
}
