# go-plotter-svg
Generate SVG files for AxiDraw plotter

`go run ./...`

See some resulting plots on [Insta](https://www.instagram.com/cube.gif/).

# TODOs

* Fix curvy line logic to work for other orientations
  * Generalize to any shape (polygon, circle, ellipse, etc.)
* Brush
	* Allow for repeating strokes every now and then
* Do a Marching Squares approach
* For truchet, try more points per edge, like 6 (laid out 1-2-1-2, which has 4 non-intersecting tiles)