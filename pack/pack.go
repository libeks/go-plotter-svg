package pack

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/libeks/go-plotter-svg/primitives"

	"github.com/kelindar/bitmap"
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
	boxes          []box // all search states will point to the same slice of boxes, they shouldn't be changing
	unprocessables bitmap.Bitmap
	positions      []primitives.Vector
	unprocessed    bitmap.Bitmap
	processed      map[int]primitives.Vector // map from box index to the translation vector
}

func (s searchState) RemoveSuperfluousPositions() searchState {
	if len(s.positions) == 0 {
		return s
	}
	newPositions := []primitives.Vector{}
	for _, pos := range s.positions {
		pt := primitives.Origin.Add(pos)
		conflict := false
		for i, v := range s.processed {
			bx := s.boxes[i].Translate(v)
			if bx.PointInside(pt) {
				// fmt.Printf("%v is inside %v\n", pt, bx)
				conflict = true
				break
			}
		}
		if !conflict {
			newPositions = append(newPositions, pos)
		} else {
			// fmt.Printf("Removing candidate position %v\n", pos)
		}
	}
	s.positions = newPositions
	return s
}

// Key is used to compare whether two states are the same
func (s searchState) Key() string {
	unprocessables, err := s.unprocessables.MarshalJSON()
	if err != nil {
		panic("Couldn't marshal unprocessables")
	}
	unprocessed, err := s.unprocessed.MarshalJSON()
	if err != nil {
		panic("Couldn't marshal unprocessed")
	}
	processedKeys := maps.Keys(s.processed)
	sort.Ints(processedKeys)
	processedStrings := make([]string, len(processedKeys))
	for i, key := range processedKeys {
		v := s.processed[key]
		processedStrings[i] = fmt.Sprintf("%d#%s", key, v.Repr())
	}
	sort.Slice(s.positions, func(i, j int) bool {
		if s.positions[i].X != s.positions[j].X {
			return s.positions[i].X > s.positions[j].X
		}
		return s.positions[i].Y != s.positions[j].Y
	})
	positionStrings := make([]string, len(s.positions))
	for i, pos := range s.positions {
		positionStrings[i] = pos.Repr()
	}
	return fmt.Sprintf("%s_%s_%s_%s", unprocessables, unprocessed, strings.Join(processedStrings, ";"), strings.Join(positionStrings, ";"))
}

// return a string that represents the bitmap of the boxes that have been placed
func (s searchState) processedKey() string {
	bitmap := bitmap.Bitmap{}
	for key := range s.processed {
		bitmap.Set(uint32(key))
	}
	key, err := bitmap.MarshalJSON()
	if err != nil {
		panic("Could not marshal")
	}
	return string(key)
}

func (s searchState) IntersectsFixed(b primitives.BBox) bool {
	for i, v := range s.processed {
		fixedBox := s.boxes[i].Translate(v).WithPadding(-199) // grow the box by the padding amount (minus a bit)
		// fmt.Printf("Box is %v, but padded is %v\n", box, fixedBox)
		if fixedBox.DoesIntersect(b) {
			// there is a conflict
			// fmt.Printf("There is an intersection between %v and %v\n", fixedBox, b)
			return true
		}
	}
	return false
}

func (s searchState) ProcessedIndexes() []int {
	ret := []int{}
	for i := range s.processed {
		ret = append(ret, i)
	}
	return ret
}

func (s searchState) UnprocessedIndexes() []int {
	ret := []int{}
	s.unprocessed.Range(func(x uint32) {
		ret = append(ret, int(x))
	})
	return ret
}

// return the area of the bounding box that contains all positioned processed boxes
func (s searchState) ProcessedArea() float64 {
	points := []primitives.Point{}
	for i, v := range s.processed {
		bx := s.boxes[i].Translate(v)
		points = append(points, bx.BBox.Corners()...)
	}
	totalBoundingBox := primitives.BBoxAroundPoints(points...)
	return totalBoundingBox.Area()
}

// return the total area of each individual placed box
func (s searchState) ProcessedAreaSum() float64 {
	total := 0.0
	for i := range s.processed {
		total += s.boxes[i].Area()
	}
	return total
}

func (s searchState) unprocessedArea() float64 {
	total := 0.0
	s.unprocessables.Range(func(id uint32) {
		total += s.boxes[id].Area()
	})
	return total
}

type PackingSolution struct {
	Translations   map[int]primitives.Vector
	DebugPositions []primitives.Vector // only for debugging, lists the positions where additional boxes could be placed
}

// given a list of bounding boxes, return a list of vectors that tell how much each box should be translated for efficient packing onto
// the container. If it cannot be placed into the container... return something more interesting
// Include a padding between the boxes
func PackOnOnePageExhaustive(bxs []primitives.BBox, container primitives.BBox, padding float64) PackingSolution {
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
	searchStates := []searchState{
		{
			boxes:          boxes,
			unprocessables: unprocessableBoxes,
			positions:      []primitives.Vector{container.UpperLeft.Subtract(primitives.Origin)},
			unprocessed:    allBoxes,
			processed:      map[int]primitives.Vector{},
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
					unprocessables: state.unprocessed.Clone(nil),
					positions:      state.positions,
					unprocessed:    bitmap.Bitmap{},
					processed:      maps.Clone(state.processed),
				})
			}
			for posID, pos := range state.positions {
				state.unprocessed.Range(func(boxID uint32) {
					boxCandidate := state.boxes[boxID]
					positionedBox := boxCandidate.Translate(pos)
					if !container.Contains(positionedBox.BBox) {
						return
					}
					if state.IntersectsFixed(positionedBox.BBox) {
						return
					}
					// box can be placed here, let's do so
					found = true
					positions := append(append([]primitives.Vector{}, state.positions[:posID]...), state.positions[posID+1:]...) // remove the current position
					// add padded positions at the two corners
					swCandidate := positionedBox.NECorner().Add(primitives.Vector{X: 0, Y: padding}).Subtract(primitives.Origin)
					neCandidate := positionedBox.SWCorner().Add(primitives.Vector{X: padding, Y: 0}).Subtract(primitives.Origin)
					positions = append(positions, neCandidate)
					if neCandidate.Y > container.UpperLeft.Y {
						candidate := primitives.Vector{X: neCandidate.X, Y: container.UpperLeft.Y}
						// fmt.Printf("Added ne candidate tangent %v from %v\n", candidate, neCandidate)
						positions = append(positions, candidate)
					}
					positions = append(positions, swCandidate)
					if swCandidate.X > container.UpperLeft.X {
						candidate := primitives.Vector{X: container.UpperLeft.X, Y: swCandidate.Y}
						// fmt.Printf("Added sw candidate tangent %v from %v\n", candidate, swCandidate)
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
					newState = newState.RemoveSuperfluousPositions()
					newStates = append(newStates, newState)
				})

			}
			if !found {
				finalStates = append(finalStates, searchState{
					boxes:          state.boxes,
					unprocessables: state.unprocessed.Clone(nil),
					positions:      append([]primitives.Vector{}, state.positions...),
					unprocessed:    bitmap.Bitmap{},
					processed:      maps.Clone(state.processed),
				})
			}
		}
		searchStates = newStates
	}
	if len(finalStates) == 0 {
		return PackingSolution{
			make(map[int]primitives.Vector, len(bxs)),
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
	fmt.Printf("Best solution has %d unplaceable, %d placeable boxes\n", len(solution.unprocessables), len(solution.processed))
	return PackingSolution{
		Translations:   solution.processed,
		DebugPositions: solution.positions,
	}
}

type areaStruct struct {
	area float64
	searchState
}

func consolidateSearchStates(states []searchState) []searchState {
	if len(states) < 2 {
		return states
	}
	fmt.Printf("Consolidating from %d... ", len(states))
	// consolidate states that contain the same boxes and cover the same area
	// stateMap := map[string]searchState{}
	// for _, state := range states {
	// 	key := fmt.Sprintf("%s:%f", state.processedKey(), state.ProcessedArea())
	// 	stateMap[key] = state
	// }
	// fmt.Printf("to %d ", len(stateMap))
	// states = maps.Values(stateMap)

	// consolidate states that leave 50%+ space empty
	newStates := []searchState{}
	for _, state := range states {
		if len(state.processed) > 3 {
			coveredArea := state.ProcessedArea()
			processedArea := state.ProcessedAreaSum()
			if processedArea/coveredArea > 0.70 {
				newStates = append(newStates, state)
			}
		} else {
			newStates = append(newStates, state)
		}
	}
	states = newStates
	fmt.Printf("to %d ", len(states))

	// // consolidate states that contain the same boxes to only keep the one that covers the least area
	// stateMap2 := map[string]areaStruct{}
	// for _, state := range states {
	// 	area := state.ProcessedArea()
	// 	key := state.processedKey()
	// 	if alternate, ok := stateMap2[key]; ok {
	// 		if alternate.area > area {
	// 			stateMap2[key] = areaStruct{
	// 				area:        area,
	// 				searchState: state,
	// 			}
	// 		}
	// 	} else {
	// 		stateMap2[key] = areaStruct{
	// 			area:        area,
	// 			searchState: state,
	// 		}
	// 	}
	// 	// stateMap[key] = state
	// }
	// states = make([]searchState, 0, len(stateMap))
	// for _, value := range stateMap2 {
	// 	states = append(states, value.searchState)
	// }
	// fmt.Printf("to %d", len(states))

	// states = maps.Values(stateMap)
	// // remove states that are identical
	// stateMap = map[string]searchState{}
	// for _, state := range states {
	// 	// fmt.Printf("key %s\n", state.Key())
	// 	stateMap[state.Key()] = state
	// }
	// fmt.Printf("to %d ", len(stateMap))
	// states = maps.Values(stateMap)

	fmt.Printf("\n")
	return states
}

// func PackOnMultiplePages(boxes []primitives.BBox, container primitives.BBox)
