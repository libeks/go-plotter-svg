package fonts

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/kintar/etxt/efixed"
	"github.com/libeks/go-plotter-svg/primitives"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	fontHeight = 100
)

type Font struct {
	*truetype.Font
}

func LoadFont(filename string) (*Font, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}
	return &Font{
		Font: font,
	}, nil
}

func convertBox(r fixed.Rectangle26_6) (primitives.Point, primitives.Point) {
	return primitives.Point{
			X: efixed.ToFloat64(r.Min.X), Y: efixed.ToFloat64(r.Min.Y),
		}, primitives.Point{
			X: efixed.ToFloat64(r.Max.X), Y: efixed.ToFloat64(r.Max.Y),
		}
}

func (f *Font) LoadGlyph(r rune) (Glyph, error) {
	index := f.Index(r)
	glyph := truetype.GlyphBuf{}
	a, _ := efixed.FromFloat64(fontHeight)
	err := glyph.Load(f.Font, a, index, font.HintingNone)
	if err != nil {
		return Glyph{}, err
	}
	// face := truetype.NewFace(f.Font, nil)
	// bounds, advance, ok := face.GlyphBounds(r)
	// if !ok {
	// 	return Glyph{}, fmt.Errorf("Couldn't get glyph bounds for rune '%v'", r)
	// }
	// minP, maxP := convertBox(glyph.Bounds)
	fmt.Printf("Glyph '%s'\n", string(r))
	// fmt.Printf("bounds %v, %v\n", minP, maxP)
	// fmt.Printf("advance %v\n", efixed.ToFloat64(glyph.AdvanceWidth))
	// fmt.Printf("Points %d, ends %d\n", len(glyph.Points), len(glyph.Ends))

	// for i, end := range glyph.Ends {
	// 	fmt.Printf("Countour %d, %v\n", i, end)
	// }

	return Glyph{
		Rune:    r,
		glyph:   glyph,
		bounds:  glyph.Bounds,
		advance: glyph.AdvanceWidth,
	}, nil
}

func (f *Font) Kerning(h float64, r1, r2 rune) float64 {
	i1 := f.Index(r1)
	i2 := f.Index(r2)
	r := h / fontHeight
	k := f.Kern(fontHeight, i1, i2)
	return r * efixed.ToFloat64(k)

}
