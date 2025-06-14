package svg

import (
	"fmt"
	"os"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/scenes"
)

type SVG struct {
	Fname  string
	Width  string
	Height string
	scenes.Scene
}

func (s SVG) WriteSVG() {
	f, err := os.OpenFile(s.Fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := xmlwriter.Open(f, xmlwriter.WithIndentString("    "))
	ec := &xmlwriter.ErrCollector{}
	defer ec.Panic()
	layers := []xmlwriter.Writable{}
	for i, layer := range s.Scene.GetLayers() {
		layers = append(layers, layer.XML(i))
	}
	ec.Do(
		w.Start(xmlwriter.Doc{}),
		w.Start(xmlwriter.Elem{
			Name: "svg", Attrs: []xmlwriter.Attr{
				{Name: "width", Value: s.Width},
				{Name: "height", Value: s.Height},
				{Name: "viewBox", Value: "0 0 10000 10000"},
				{Name: "version", Value: "1.1"},
				{Name: "id", Value: "svg6"},
				{Name: "sodipodi:docname", Value: "test_inkscape.svg"},
				{Name: "inkscape:version", Value: "1.3.2 (091e20e, 2023-11-25, custom)"},
				{Name: "xmlns:inkscape", Value: "http://www.inkscape.org/namespaces/inkscape"},
				{Name: "xmlns:sodipodi", Value: "http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"},
				{Name: "xmlns", Value: "http://www.w3.org/2000/svg"},
				{Name: "xmlns:svg", Value: "http://www.w3.org/2000/svg"},
			},
			Content: layers,
		}),
		w.EndAllFlush(),
	)
	fmt.Printf("Finished rendering to file %s\n", s.Fname)
}
