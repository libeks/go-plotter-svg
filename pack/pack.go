package pack

import (
	"fmt"
	"slices"

	"golang.org/x/exp/maps"

	"github.com/libeks/go-plotter-svg/primitives"

	"github.com/kelindar/bitmap"
)

const MAX_STATES = 1_000_000

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

// Paged Vector describes a vector that is placed on a specific page of the document
type PagedVector struct {
	primitives.Vector
	IsUpperLeftCorner bool
	Page              int
}

func (p PagedVector) Repr() string {
	return fmt.Sprintf("%d=%s", p.Page, p.Vector.Repr())
}

type PackingSolution struct {
	Pages          int
	Translations   map[int]PagedVector
	DebugPositions []PagedVector // only for debugging, lists the positions where additional boxes could be placed
}

// given a list of bounding boxes, return a list of vectors that tell how much each box should be translated for efficient packing onto
// the container. If it cannot be placed into the container... return something more interesting
// Include a padding between the boxes
func PackOnOnePage(bxs []primitives.BBox, container primitives.BBox, padding float64) PackingSolution {
	return packingAlgo(bxs, container, padding, false)
}

func PackOnMultiplePages(bxs []primitives.BBox, container primitives.BBox, padding float64) PackingSolution {
	return packingAlgo(bxs, container, padding, true)
}

// given a list of bounding boxes, return a list of vectors that tell how much each box should be translated for efficient packing onto
// the container. If it cannot be placed into the container... return something more interesting
// Include a padding between the boxes
func packingAlgo(bxs []primitives.BBox, container primitives.BBox, padding float64, multiPage bool) PackingSolution {
	// unprocessables := []box{}
	boxes := make([]box, len(bxs))

	allBoxes := bitmap.Bitmap{}
	unprocessableBoxes := bitmap.Bitmap{}
	for i, b := range bxs {
		// ensure that all boxes have their upper left at origin for easier handling
		boxes[i] = box{
			BBox:  b.Translate(b.UpperLeft.Subtract(primitives.Origin)),
			index: i,
		}
		if b.Width() > container.Width() || b.Height() > container.Height() {
			fmt.Printf("Box %v is too big for the container, will ignore\n", b)
			unprocessableBoxes.Set(uint32(i))
			continue
		}
		allBoxes.Set(uint32(i))
	}
	searchStates := []*searchState{
		{
			boxes:          boxes,
			unprocessables: unprocessableBoxes,
			positions: []PagedVector{{
				Page:              0,
				IsUpperLeftCorner: true,
				Vector:            container.UpperLeft.Subtract(primitives.Origin)}},
			unprocessed: allBoxes,
			processed:   map[int]PagedVector{},
		},
	}
	finalStates := []*searchState{}

	pass := 1
	for len(searchStates) > 0 {
		newStates := []*searchState{} // append to this slice, then swap this out with searchStates at the end
		fmt.Printf("Pass %d: have %d search states\n", pass, len(searchStates))
		if len(searchStates) > MAX_STATES {
			panic("Exceeded max allowed states, aborting")
		}
		searchStates = consolidateSearchStates(searchStates, container)
		slices.SortFunc(searchStates, stateComparatorFunc) // sort slices so this pass starts with the best options so far
		for _, state := range searchStates {
			// fmt.Printf("%v, %v\n", state.processed, state.positions)
			found := false
			if state.unprocessed.Count() == 0 {
				// if there are no unprocessed boxes, we are done with this state
				finalStates = append(finalStates, state)
				continue
			}
			if len(state.positions) == 0 {
				// if there are no positions, we are done, but the rest of the boxes are unprocessable
				finalStates = append(finalStates, &searchState{
					unprocessables: state.unprocessed.Clone(nil),
					positions:      state.positions,
					unprocessed:    bitmap.Bitmap{},
					processed:      maps.Clone(state.processed),
				})
			}
			for posID, pos := range state.positions {
				state.unprocessed.Range(func(boxID uint32) {
					boxCandidate := state.boxes[boxID]
					positionedBox := boxCandidate.Translate(pos.Vector)
					if !container.Contains(positionedBox.BBox) {
						return
					}
					if state.IntersectsFixed(pos.Page, positionedBox.BBox) {
						return
					}
					// box can be placed here, let's do so
					found = true
					positions := append(
						append(
							[]PagedVector{},
							state.positions[:posID]...,
						),
						state.positions[posID+1:]...,
					) // remove the current position
					if pos.IsUpperLeftCorner && multiPage {
						// Add a new page to positions
						positions = append(
							positions,
							PagedVector{
								Page:              pos.Page + 1,
								IsUpperLeftCorner: true,
								Vector:            container.UpperLeft.Subtract(primitives.Origin),
							},
						)
					}
					// add padded positions at the two corners
					swCandidate := PagedVector{
						Page:              pos.Page,
						IsUpperLeftCorner: false,
						Vector:            positionedBox.NECorner().Add(primitives.Vector{X: 0, Y: padding}).Subtract(primitives.Origin),
					}
					neCandidate := PagedVector{
						Page:              pos.Page,
						IsUpperLeftCorner: false,
						Vector:            positionedBox.SWCorner().Add(primitives.Vector{X: padding, Y: 0}).Subtract(primitives.Origin),
					}
					positions = append(positions, neCandidate)
					if neCandidate.Y > container.UpperLeft.Y {
						candidate := PagedVector{
							Page:              pos.Page,
							IsUpperLeftCorner: false,
							Vector:            primitives.Vector{X: neCandidate.X, Y: container.UpperLeft.Y},
						}
						positions = append(positions, candidate)
					}
					positions = append(positions, swCandidate)
					if swCandidate.X > container.UpperLeft.X {
						candidate := PagedVector{
							Page:              pos.Page,
							IsUpperLeftCorner: false,
							Vector:            primitives.Vector{X: container.UpperLeft.X, Y: swCandidate.Y},
						}
						positions = append(positions, candidate)
					}

					newProcessed := maps.Clone(state.processed)
					newProcessed[int(boxID)] = pos
					newUnprocessed := state.unprocessed.Clone(nil)
					newUnprocessed.Remove(boxID)
					newState := searchState{
						boxes:          state.boxes,
						unprocessables: state.unprocessables.Clone(nil),
						positions:      positions,
						unprocessed:    newUnprocessed, // strip out this box
						processed:      newProcessed,
					}
					newState = newState.RemoveSuperfluousPositions(container)
					newStates = append(newStates, &newState)
				})
			}
			if !found {
				finalStates = append(finalStates, &searchState{
					boxes:          state.boxes,
					unprocessables: state.unprocessed.Clone(nil),
					positions:      append([]PagedVector{}, state.positions...),
					unprocessed:    bitmap.Bitmap{},
					processed:      maps.Clone(state.processed),
				})
			}
			if len(newStates) > MAX_STATES {
				fmt.Printf("Stopping pass early since we already have %d potential cases\n", len(newStates))
				newStates = newStates[:MAX_STATES]
				break // stop adding more candidates if max states is exceeded
			}
		}
		searchStates = newStates
		pass += 1
	}
	if len(finalStates) == 0 {
		return PackingSolution{
			Pages:          0,
			Translations:   make(map[int]PagedVector, len(bxs)),
			DebugPositions: []PagedVector{},
		}
	}
	fmt.Printf("Got %d possible solutions\n", len(finalStates))
	solution := slices.MinFunc(finalStates, stateComparatorFunc)
	fmt.Printf("Best solution has %d unplaceable, %d placeable boxes\n", solution.unprocessables.Count(), len(solution.processed))

	fmt.Printf("Boxes are:\n")
	for i, v := range solution.processed {
		fmt.Printf("\t%d: %v\n", v.Page, boxes[i].Translate(v.Vector))
	}
	fmt.Printf("\nPositions are:\n")
	for _, pos := range solution.positions {
		fmt.Printf("\t%d: %v\n", pos.Page, pos.Vector)
	}
	return PackingSolution{
		Pages:          solution.Pages(),
		Translations:   solution.processed,
		DebugPositions: solution.positions,
	}
}

// func PackOnMultiplePages(boxes []primitives.BBox, container primitives.BBox)
