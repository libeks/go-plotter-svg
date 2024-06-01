# go-plotter-svg
Generate SVG files for AxiDraw plotter

`go run ./...`

See some resulting plots on [Insta](https://www.instagram.com/cube.gif/).

# TODOs

* Fix curvy line logic to work for other orientations
  * Generalize to any shape (polygon, circle, ellipse, etc.)
* Add Test Page to calibrate different pen alignment
  * It could also be an optional print on the edge of the page
* Add ability to configure different pen alignment
* Brush
	* Allow for repeating strokes every now and then
* Rethink layer abstraction, it doesn't really work now
* Write own svg marshaller, svgo only works in integer space, whereas svg does support floating point
