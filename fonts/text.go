package fonts

import (
	"fmt"
	"runtime"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

type option struct {
	fontFile string
	size     float64 // vertical height to render to, somewhat approximate
	// The following are currently unimplemented
	vAlignment    string  // how the text is positioned inside the bounding box vertically, default to center
	hAlignment    string  // how the text is positioned inside the bounding box horizontally, default to center
	rotationAngle float64 // clock-wise, how much the text should be rotated around its center before positioning, in radians
}

// TextRenderer contains everything that you'd need to render text to a Scene
type TextRender struct {
	Text       string
	CharBoxes  []box.Box
	CharCurves []lines.LineLike
	CharPoints []ControlPoint
	// The following is currently unimplemented
	Kernings []float64 // the amount of kerning between every pair of characters
}

func (t TextRender) Translate(v primitives.Vector) TextRender {
	newBoxes := make([]box.Box, len(t.CharBoxes))
	for i, b := range t.CharBoxes {
		newBoxes[i] = b.Translate(v)
	}
	newCurves := make([]lines.LineLike, len(t.CharCurves))
	for i, c := range t.CharCurves {
		newCurves[i] = c.Translate(v)
	}
	newPoints := make([]ControlPoint, len(t.CharPoints))
	for i, c := range t.CharPoints {
		newPoints[i] = c.Translate(v)
	}
	return TextRender{
		Text:       t.Text,
		CharBoxes:  newBoxes,
		CharCurves: newCurves,
		CharPoints: newPoints,
		// The following are currently unimplemented
		Kernings: t.Kernings,
	}
}

func (t TextRender) BoundingBox() box.Box {
	if len(t.CharBoxes) == 0 {
		return box.Box{}
	}
	b := t.CharBoxes[0].BBox()
	for _, newB := range t.CharBoxes[1:] {
		b = b.Add(newB.BBox())
	}
	return box.Box{
		X:    b.UpperLeft.X,
		Y:    b.UpperLeft.Y,
		XEnd: b.LowerRight.X,
		YEnd: b.LowerRight.Y,
	}
}

type textOption func(option) option

func WithFont(path string) textOption {
	return func(o option) option {
		o.fontFile = path
		return o
	}
}

func WithSize(size float64) textOption {
	return func(o option) option {
		o.size = size
		return o
	}
}

// RenderText renders the specified text, with the specified options, in the middle of the bounding box
// The text may span outside of the bounding box, i.e. RenderText currently doesn't resize to fit the box
func RenderText(b box.Box, text string, textOptions ...textOption) TextRender {
	var fontFile string
	if runtime.GOOS == "windows" {
		fontFile = "C:/Windows/Fonts/seguihis.ttf"
		// fontFile="C:/Windows/Fonts/bahnschrift.ttf"
	} else {
		// MacOS
		fontFile = "/System/Library/Fonts/Palatino.ttc"
	}
	o := option{
		fontFile: fontFile,
		size:     1000.0,
	}
	for _, opt := range textOptions {
		o = opt(o)
	}
	f, err := LoadFont(o.fontFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded font %s\n", f.Name(3)) // print "Unique subfamily identification"
	// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6name.html

	controlPoints := []ControlPoint{}
	curves := []lines.LineLike{}
	charBoxes := []box.Box{}
	offsetX := 0.0
	prevCh := ' '
	for _, ch := range text {
		fmt.Printf("offset %f\n", offsetX)
		glyph, err := f.LoadGlyph(ch)
		if err != nil {
			panic(err)
		}

		c := glyph.GetHeightCurves(o.size)
		c = c.Translate(primitives.Vector{X: offsetX, Y: 0})
		offsetX += c.AdvanceWidth
		fmt.Printf("Advacned Width %f\n", c.AdvanceWidth)
		kern := f.Kerning(6000, prevCh, ch)
		fmt.Printf("kerning %f\n", kern)
		offsetX += kern
		curves = append(curves, c.Curves...)
		controlPoints = append(controlPoints, c.Points...)
		charBoxes = append(charBoxes, c.BoundingBox)
		prevCh = ch
	}

	render := TextRender{
		Text:       text,
		CharCurves: curves,
		CharPoints: controlPoints,
		CharBoxes:  charBoxes,
	}
	v := b.Center().Subtract(render.BoundingBox().Center()) // center within the bounding box
	return render.Translate(v)
}
