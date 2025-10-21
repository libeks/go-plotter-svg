package svg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/scenes"
)

type SVG struct {
	Fname  string
	Width  string
	Height string
	scenes.Document
}

func (s SVG) WritePage(fname string, scene scenes.Page) {
	f, err := os.OpenFile(s.Fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := xmlwriter.Open(f, xmlwriter.WithIndentString("    "))
	ec := &xmlwriter.ErrCollector{}
	defer ec.Panic()
	layers := []xmlwriter.Writable{}
	for i, layer := range scene.GetLayers() {
		layers = append(layers, layer.XML(i))
	}
	ec.Do(
		w.Start(xmlwriter.Doc{}),
		w.Start(xmlwriter.Elem{
			Name: "svg", Attrs: []xmlwriter.Attr{
				{Name: "width", Value: s.Width},
				{Name: "height", Value: s.Height},
				{Name: "viewBox", Value: "0 0 13333 10000"}, // ensure the viewbox fits the page size
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
	fmt.Printf("Finished rendering to file %s\n", fname)
}

func (s SVG) WriteSVG() {
	if s.Document.NumPages() == 1 {
		s.WritePage(s.Fname, s.Document.Page(0))
	} else {
		extension := filepath.Ext(s.Fname)
		basename := strings.TrimSuffix(s.Fname, extension)
		fnamePattern := fmt.Sprintf("%s_%%d%s", basename, extension)
		for i := range s.Document.NumPages() {
			s.WritePage(fmt.Sprintf(fnamePattern, i), s.Document.Page(i))
		}
	}
}
