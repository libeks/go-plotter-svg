package fonts

import (
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/kintar/etxt/efixed"
	"golang.org/x/image/font"
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

func (f *Font) LoadGlyph(r rune) (Glyph, error) {
	index := f.Index(r)
	glyph := truetype.GlyphBuf{}
	a, _ := efixed.FromFloat64(fontHeight)
	err := glyph.Load(f.Font, a, index, font.HintingNone)
	if err != nil {
		return Glyph{}, err
	}

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
