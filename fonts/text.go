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
	fitToBox bool    // fit the bounding box of the text inside of the provided box, this conflicts with the `size` parameter

	// The following are currently unimplemented
	// vAlignment    string  // how the text is positioned inside the bounding box vertically, default to center
	// hAlignment    string  // how the text is positioned inside the bounding box horizontally, default to center
	// rotationAngle float64 // clock-wise, how much the text should be rotated around its center before positioning, in radians
}

// TextRenderer contains everything that you'd need to render text to a Scene
type TextRender struct {
	Text       string
	Size       float64 // size at which text was rendered
	CharBoxes  []box.Box
	CharCurves []lines.LineLike
	CharPoints []ControlPoint
	// The following is currently unimplemented
	// Kernings []float64 // the amount of kerning between every pair of characters
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
		Size:       t.Size,
		// The following are currently unimplemented
		// Kernings: t.Kernings,
	}
}

func (t TextRender) BoundingBox() box.Box {
	if len(t.CharBoxes) == 0 {
		return box.Box{}
	}
	b := t.CharBoxes[0].BBox
	for _, newB := range t.CharBoxes[1:] {
		b = b.Add(newB.BBox)
	}
	return box.Box{BBox: b}
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

func WithFitToBox() textOption {
	return func(o option) option {
		o.fitToBox = true
		return o
	}
}

// RenderText renders the specified text, with the specified options, in the middle of the bounding box
// The text may span outside of the bounding box, unless WithFitToBox is used
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

	size := o.size
	for size > 1.0 { // loop until the text is inside the box or the size is too small to render
		controlPoints := []ControlPoint{}
		curves := []lines.LineLike{}
		charBoxes := []box.Box{}
		offsetX := 0.0
		prevCh := ' '
		for _, ch := range text {
			glyph, err := f.LoadGlyph(ch)
			if err != nil {
				panic(err)
			}

			c := glyph.GetHeightCurves(size)
			c = c.Translate(primitives.Vector{X: offsetX, Y: 0})
			offsetX += c.AdvanceWidth
			kern := f.Kerning(6000, prevCh, ch)
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
			Size:       size,
		}
		v := b.Center().Subtract(render.BoundingBox().Center()) // center within the bounding box
		textBox := render.Translate(v)

		boundingBox := textBox.BoundingBox()
		hRatio := boundingBox.Width() / b.Width()
		vRatio := boundingBox.Height() / b.Height()
		if !o.fitToBox || (hRatio <= 1.0 && vRatio <= 1.0) {
			// if need to fit the box and the character is bigger than the box, scale down, otherwise return
			fmt.Printf("Rendered at size %f\n", size)
			return textBox
		}
		size = size / max(hRatio, vRatio)
	}
	fmt.Printf("Box is too small to render text\n")
	return TextRender{}
}
