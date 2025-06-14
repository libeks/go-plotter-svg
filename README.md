# go-plotter-svg
Generate SVG files for AxiDraw plotter

`go run ./...`

See some resulting plots on [Insta](https://www.instagram.com/cube.gif/).

# TODOs

* Fix curvy line logic to work for other orientations
  * Generalize to any shape (polygon, circle, ellipse, etc.)
* Brush
	* Allow for repeating strokes every now and then
* Come up with a solution for object-to-curve intersection
  * The current appraoch only works for lines and circle segments, but does not for Beziers. To generalize, I'll have to step through each line iteratively, though the question is, how densely should I sample for each object? Is there some easy way to find out the "depth" of a point, both inside and outside?
* Investigate whether stroke speed has an impact on pen performance
* Implement closed-curve startpoint shuffling
  * Requires a way to arbitrarily subdivide a linelike
* Use Marching Squares off of an object-distance measure field
* make sense of what counts as clockwise w.r.t. circle arc angles. Is the angle measures CCW? it doesn't make sense now. Maybe vectors.RotateCCW is wrong?
* Font rendering - distinguish inside vs outside of a glyph contour using winding number. This is relevant for bandshift 'a', which has overlapping contours, vs 'o'
  https://developer.apple.com/fonts/TrueType-Reference-Manual/RM02/Chap2.html#distinguishing
* For font families, consider this library:
  https://pkg.go.dev/github.com/benoitkugler/go-opentype#section-readme
* RenderText - add proper spacing between characters

* Foldables:
  * Refactor to be easier to specify
    * Something like: Here are the numbered faces, here are where they are connected with each other, here are where we have flaps, etc
  * Add rendering of face numbers
  * Properly render the rectangle around what can be plotted

# Plot ideas

* Halftone circles in different colors
* Marching squares
* Concentric circles around several points, all clipped to the respective Voronoi diagram
