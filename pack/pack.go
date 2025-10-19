package pack

import (
	"cmp"
	"fmt"
	"slices"
	"sort"

	"github.com/libeks/go-plotter-svg/primitives"
)

type box struct {
	primitives.BBox
	index int
}

func (b box) Translate(v primitives.Vector) box {
	return box{
		BBox:  b.BBox.Translate(v),
		index: b.index,
	}
}

type fixedBox struct {
	boxes []primitives.BBox
}

func (fb fixedBox) Intersect(b primitives.BBox) bool {
	for _, fixedBox := range fb.boxes {
		if fixedBox.DoesIntersect(b) {
			// there is a conflict
			// fmt.Printf("Box %v intersects with box %v\n", fixedBox, b)
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
				if !container.Contains(positionedBox.BBox) {
					fmt.Printf("positioned box %v is not inside the container, skipping\n", positionedBox)
					continue
				}
				if fixedBoxes.Intersect(positionedBox.BBox) {
					fmt.Printf("positioned box %v conflicts with an existing box, skipping\n", positionedBox)
					continue
				}
				// box can be placed here, let's do so
				found = true
				foundOuter = true
				fmt.Printf("Found a solution %v\n", pos)
				fixedBoxes.boxes = append(fixedBoxes.boxes, positionedBox.BBox)

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

type searchState struct {
	unprocessables []box
	positions      []primitives.Vector
	unprocessed    []box
	processed      []box
}

func (s searchState) RemoveSuperfluousPositions() searchState {
	if len(s.positions) == 0 {
		return s
	}
	newPositions := []primitives.Vector{}
	for _, pos := range s.positions {
		pt := primitives.Origin.Add(pos)
		conflict := false
		for _, bx := range s.processed {
			if bx.PointInside(pt) {
				conflict = true
				break
			}
		}
		if !conflict {
			newPositions = append(newPositions, pos)
		} else {
			// fmt.Printf("Found superfluous position\n")
		}
	}
	s.positions = newPositions
	return s
}

// // Key is used to compare whether two states are the same
// func (s searchState) Key() string {

// }

func (s searchState) Copy() searchState {
	return searchState{
		unprocessables: append([]box{}, s.unprocessables...),
		positions:      append([]primitives.Vector{}, s.positions...),
		unprocessed:    append([]box{}, s.unprocessed...),
		processed:      append([]box{}, s.processed...),
	}
}

func (s searchState) IntersectsFixed(b primitives.BBox) bool {
	for _, fixedBox := range s.processed {
		if fixedBox.DoesIntersect(b) {
			// there is a conflict
			// fmt.Printf("Box %v intersects with box %v\n", fixedBox, b)
			return true
		}
	}
	return false
}

func (s searchState) ProcessedIndexes() []int {
	ret := []int{}
	for _, box := range s.processed {
		ret = append(ret, box.index)
	}
	return ret
}

func (s searchState) UnprocessedIndexes() []int {
	ret := []int{}
	for _, box := range s.unprocessed {
		ret = append(ret, box.index)
	}
	return ret
}

// return the area of the bounding box that contains all processed boxes
func (s searchState) ProcessedArea() float64 {
	points := []primitives.Point{}
	for _, bx := range s.processed {
		points = append(points, bx.BBox.Corners()...)
	}
	totalBoundingBox := primitives.BBoxAroundPoints(points...)
	return totalBoundingBox.Area()
}

func (s searchState) unprocessedArea() float64 {
	total := 0.0
	for _, box := range s.unprocessables {
		total += box.Area()
	}
	return total
}

type PackingSolution struct {
	Translations   []*primitives.Vector
	DebugPositions []primitives.Vector // only for debugging, lists the positions where additional boxes could be placed
}

// given a list of bounding boxes, return a list of vectors that tell how much each box should be translated for efficient packing onto
// the container. If it cannot be placed into the container... return something more interesting
// Include a padding between the boxes
func PackOnOnePageExhaustive(bxs []primitives.BBox, container primitives.BBox, padding float64) PackingSolution {
	unprocessables := []box{}
	boxes := make([]box, len(bxs))
	for i, b := range bxs {
		if b.Width() > container.Width() || b.Height() > container.Height() {
			fmt.Printf("Box %v is too big for the container, will ignore\n", b)
			unprocessables = append(unprocessables,
				box{
					BBox:  b.Translate(b.UpperLeft.Subtract(primitives.Origin)),
					index: i,
				},
			)
			continue
		}
		// ensure that all boxes have their upper left at origin for easier handling
		boxes[i] = box{
			BBox:  b.Translate(b.UpperLeft.Subtract(primitives.Origin)),
			index: i,
		}
	}
	// fmt.Printf("boxes %v\n", boxes)
	searchStates := []searchState{
		{
			unprocessables: unprocessables,
			positions:      []primitives.Vector{container.UpperLeft.Subtract(primitives.Origin)},
			unprocessed:    boxes,
			processed:      []box{},
		},
	}
	finalStates := []searchState{}
	for len(searchStates) > 0 {
		newStates := []searchState{} // append to this slice, then swap this out with searchStates at the end
		fmt.Printf("have %d search states\n", len(searchStates))
		searchStates = consolidateSearchStates(searchStates)
		for _, state := range searchStates {
			found := false
			if len(state.unprocessed) == 0 {
				// if there are no unprocessed boxes, we are done with this state
				finalStates = append(finalStates, state)
				continue
			}
			if len(state.positions) == 0 {
				// if there are no positions, we are done, but the rest of the boxes are unprocessable
				finalStates = append(finalStates, searchState{
					unprocessables: append([]box{}, state.unprocessed...),
					positions:      state.positions,
					unprocessed:    []box{},
					processed:      append([]box{}, state.processed...),
				})
			}
			for posID, pos := range state.positions {
				for boxID, boxCandidate := range state.unprocessed {
					positionedBox := boxCandidate.Translate(pos)
					if !container.Contains(positionedBox.BBox) {
						continue
					}
					if state.IntersectsFixed(positionedBox.BBox) {
						continue
					}
					// box can be placed here, let's do so
					found = true
					positions := append(append([]primitives.Vector{}, state.positions[:posID]...), state.positions[posID+1:]...) // remove the current position
					// add padded positions at the two corners
					positions = append(positions, positionedBox.NECorner().Add(primitives.Vector{X: 0, Y: padding}).Subtract(primitives.Origin))
					positions = append(positions, positionedBox.SWCorner().Add(primitives.Vector{X: padding, Y: 0}).Subtract(primitives.Origin))

					newState := searchState{
						unprocessables: append([]box{}, state.unprocessables...),
						positions:      positions,
						unprocessed:    append(append([]box{}, state.unprocessed[:boxID]...), state.unprocessed[boxID+1:]...), // strip out this box
						processed:      append(append([]box{}, state.processed...), positionedBox),
					}
					newState = newState.RemoveSuperfluousPositions()
					newStates = append(newStates, newState)
				}

			}
			if !found {
				finalStates = append(finalStates, searchState{
					unprocessables: append([]box{}, state.unprocessed...),
					positions:      append([]primitives.Vector{}, state.positions...),
					unprocessed:    []box{},
					processed:      append([]box{}, state.processed...),
				})
			}
		}
		searchStates = newStates
	}
	if len(finalStates) == 0 {
		return PackingSolution{
			make([]*primitives.Vector, len(bxs)),
			[]primitives.Vector{},
		}
	}
	fmt.Printf("Got %d possible solutions\n", len(finalStates))
	solution := slices.MinFunc(finalStates, func(a, b searchState) int {
		// there should be as little unprocessable
		if len(a.unprocessables) != len(b.unprocessables) {
			return cmp.Compare(a.unprocessedArea(), b.unprocessedArea())
		}
		// sort by total area of processed
		return cmp.Compare(a.ProcessedArea(), b.ProcessedArea())
	})
	// for i, solution := range finalStates {
	// 	fmt.Printf("%d, %d unprocessables, area: %f\n", i, len(solution.unprocessables), solution.ProcessedArea())
	// }
	// solution := finalStates[0]
	// fmt.Printf("Best solution is %v\n", solution)
	fmt.Printf("Best solution has %d unplaceable, %d placeable boxes\n", len(solution.unprocessables), len(solution.processed))
	// fmt.Printf("%v\n", solution.processed)
	positions := make([]*primitives.Vector, len(bxs))
	for _, bx := range solution.processed {
		tmp := bx.UpperLeft.Subtract(bxs[bx.index].UpperLeft)
		// fmt.Printf("processed box %d (%d) has tranlsation of %v\n", i, bx.index, tmp)
		positions[bx.index] = &tmp
	}
	return PackingSolution{
		Translations:   positions,
		DebugPositions: solution.positions,
	}
}

func consolidateSearchStates(states []searchState) []searchState {
	if len(states) < 2 {
		return states
	}
	return states
}

// func PackOnMultiplePages(boxes []primitives.BBox, container primitives.BBox)
