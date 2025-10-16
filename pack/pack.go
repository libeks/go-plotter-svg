package pack

import (
	"github.com/libeks/go-plotter-svg/primitives"
)

// given a list of bounding boxes, return a list of vectors that tell how much each box should be translated for efficient packing onto
// the container. If it cannot be placed into the container... return something more interesting
func PackOnOnePage(boxes []primitives.BBox, container primitives.BBox) []primitives.Vector {
	vects := make([]primitives.Vector, len(boxes))
	return vects
}

// func PackOnMultiplePages(boxes []primitives.BBox, container primitives.BBox)
