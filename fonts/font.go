package fonts

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/kintar/etxt/efixed"
	"golang.org/x/image/font"
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
	a, _ := efixed.FromFloat64(6000)
	err := glyph.Load(f.Font, a, index, font.HintingNone)
	if err != nil {
		return Glyph{}, err
	}
	fmt.Printf("Glyph %s\n", string(r))
	fmt.Printf("bounds %v\n", glyph.Bounds)
	fmt.Printf("Points %d, ends %d\n", len(glyph.Points), len(glyph.Ends))

	for i, end := range glyph.Ends {
		fmt.Printf("Countour %d, %v\n", i, end)

	}

	return Glyph{
		glyph: glyph,
	}, nil
}
