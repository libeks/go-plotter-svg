# go-plotter-svg
Generate SVG files for AxiDraw plotter

`go run ./...`

See the generated outputs in [the gallery](https://github.com/libeks/go-plotter-svg/tree/main/gallery).
See some resulting plots on [Insta](https://www.instagram.com/cube.gif/).

# TODOs

* Fix curvy line logic to work for other orientations
  * Generalize to any shape (polygon, circle, ellipse, etc.)
* Brush
	* Allow for repeating strokes every now and then
* Come up with a solution for object-to-curve intersection
  * The current appraoch only works for lines and circle segments, but does not for Beziers. To generalize, I'll have to step through each line iteratively, though the question is, how densely should I sample for each object? Is there some easy way to find out the "depth" of a point, both inside and outside?
* Investigate whether stroke speed has an impact on pen performance
* Use Marching Squares off of an object-distance measure field
* make sense of what counts as clockwise w.r.t. circle arc angles. Is the angle measures CCW? it doesn't make sense now. Maybe vectors.RotateCCW is wrong?
* Font rendering - distinguish inside vs outside of a glyph contour using winding number. This is relevant for bandshift 'a', which has overlapping contours, vs 'o'
  https://developer.apple.com/fonts/TrueType-Reference-Manual/RM02/Chap2.html#distinguishing
* For font families, consider this library:
  https://pkg.go.dev/github.com/benoitkugler/go-opentype#section-readme
* Foldable:
  * Detect face overlap, such as with flaps being too wide, etc.
  * Allow connections between foldable objects, i.e. have a scene contain multiple foldables with interlinking
  * See what damage comes from mixing clockwise- and counter-clockwise faces
* Rectangle-packing: Optimized coveraged-based filter to preserve 10% no matter what
* Curve optimization - don't apply to fill layers, when the ink dries and is covered yet again, it becomes darker. All fill should be applied in order, preferably with alternating directions
* Allow not rendering guide for certain layers

# Plot ideas

* Halftone circles in different colors
* Marching squares
* Concentric circles around several points, all clipped to the respective Voronoi diagram
* A plane tiled into polygons, glued together with flaps