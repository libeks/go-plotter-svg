package pack

import (
	"fmt"
	"sort"

	"github.com/libeks/go-plotter-svg/primitives"
)

type box struct {
	primitives.BBox
	index int
}

type fixedBox struct {
	boxes []primitives.BBox
}

func (fb fixedBox) Intersect(b primitives.BBox) bool {
	for _, fixedBox := range fb.boxes {
		if fixedBox.DoesIntersect(b) {
			// there is a conflict
			fmt.Printf("Box %v intersects with box %v\n", fixedBox, b)
			return true
		}
	}
	return false
}

// given a list of bounding boxes, return a list of vectors that tell how much each box should be translated for efficient packing onto
// the container. If it cannot be placed into the container... return something more interesting
// Include a padding between the boxes
func PackOnOnePage(bxs []primitives.BBox, container primitives.BBox, padding float64) []*primitives.Vector {
	boxes := make([]box, len(bxs))
	for i, b := range bxs {
		if b.Width() > container.Width() || b.Height() > container.Height() {
			fmt.Printf("Box %v is too big for the container, will ignore\n", b)
			continue
		}
		// ensure that all boxes have their upper left at origin for easier handling
		boxes[i] = box{
			BBox:  b.Translate(b.UpperLeft.Subtract(primitives.Origin)),
			index: i,
		}
	}
	// sort by width
	sort.Slice(boxes, func(i, j int) bool {
		if boxes[i].Width() == boxes[j].Width() {
			return boxes[i].Height() > boxes[j].Height()
		}
		return boxes[i].Width() > boxes[j].Width()
	})

	for i, box := range boxes {
		fmt.Printf("%d Box %v\n", i, box)
	}
	vects := make([]*primitives.Vector, len(bxs))
	positions := []primitives.Vector{container.UpperLeft.Subtract(primitives.Origin)}
	fixedBoxes := fixedBox{[]primitives.BBox{}}
	for len(boxes) > 0 && len(positions) > 0 {
		foundOuter := false
		fmt.Printf("Positions %v\n", positions)
		for posID, pos := range positions {
			found := false
			for boxID, box := range boxes {
				fmt.Printf("posID %d, boxID %d (%d, %d)\n", posID, boxID, len(positions), len(boxes))
				positionedBox := box.Translate(pos)
				if !container.Contains(positionedBox) {
					fmt.Printf("positioned box %v is not inside the container, skipping\n", positionedBox)
					continue
				}
				if fixedBoxes.Intersect(positionedBox) {
					fmt.Printf("positioned box %v conflicts with an existing box, skipping\n", positionedBox)
					continue
				}
				// box can be placed here, let's do so
				found = true
				foundOuter = true
				fmt.Printf("Found a solution %v\n", pos)
				fixedBoxes.boxes = append(fixedBoxes.boxes, positionedBox)

				vects[box.index] = &pos
				positions = append(positions[:posID], positions[posID+1:]...) // remove the current position
				// add padded positions at the two corners
				positions = append(positions, positionedBox.NECorner().Add(primitives.Vector{X: 0, Y: padding}).Subtract(primitives.Origin))
				positions = append(positions, positionedBox.SWCorner().Add(primitives.Vector{X: padding, Y: 0}).Subtract(primitives.Origin))

				boxes = append(boxes[:boxID], boxes[boxID+1:]...)
				break
			}
			if found {
				break
			}
		}
		if !foundOuter {
			panic("Couldn't find a soution")
		}
	}

	return vects
}

// func PackOnMultiplePages(boxes []primitives.BBox, container primitives.BBox)
