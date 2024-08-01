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
* Come up with a solution for object-to-curve intersection
  * The current appraoch only works for lines and circle segments, but does not for Beziers. To generalize, I'll have to step through each line iteratively, though the question is, how densely should I sample for each object? Is there some easy way to find out the "depth" of a point, both inside and outside?
* Allow for curve rearranging, including reversing a curve
* Fix circle arc abstraction, move to take in circle and t-values instead of the current implementation

# Plot ideas

* Halftone circles in different colors
* Marching squares
* Concentric circles around several points, all clipped to the respective Voronoi diagram