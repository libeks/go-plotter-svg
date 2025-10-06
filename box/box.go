package box

import (
	"github.com/libeks/go-plotter-svg/primitives"
)

// returns the relative [-1,1] range in which this box's center exists in the parent box
func RelativeMinusPlusOneCenter(b, parentBox primitives.BBox) primitives.Point {
	center := b.Center()
	parentCenter := parentBox.Center()
	return primitives.Point{
		X: 2 * (center.X - parentCenter.X) / parentBox.Width(),
		Y: 2 * (center.Y - parentCenter.Y) / parentBox.Height(),
	}
}
