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

// Paged Vector describes a vector that is placed on a specific page of the document
type PagedVector struct {
	primitives.Vector
	IsUpperLeftCorner bool
	Page              int
}

func (p PagedVector) Repr() string {
	return fmt.Sprintf("%d=%s", p.Page, p.Vector.Repr())
}

type searchState struct {
	boxes          []box // all search states will point to the same slice of boxes, they shouldn't be changing
	_pages         int   // cached value of number of pages, if 0 it is not cached
	unprocessables bitmap.Bitmap
	positions      []PagedVector
	unprocessed    bitmap.Bitmap
	processed      map[int]PagedVector // map from box index to the translation vector

}

func (s *searchState) Pages() int {
	// check if value is cached
	if s._pages > 0 {
		// fmt.Printf("from cache\n")
		return s._pages
	}
	// fmt.Printf("recomputing\n")
	pages := 1
	for _, processed := range s.processed {
		if processed.Page+1 > pages {
			pages = processed.Page + 1
		}
	}
	s._pages = pages
	return pages
}

func (s searchState) RemoveSuperfluousPositions() searchState {
	if len(s.positions) == 0 {
		return s
	}
	newPositions := []PagedVector{}
	for _, pos := range s.positions {
		pt := primitives.Origin.Add(pos.Vector)
		conflict := false
		for i, v := range s.processed {
			if v.Page != pos.Page { // if the box is not on the same page as the position, ignore; there is no conflict
				continue
			}
			bx := s.boxes[i].Translate(v.Vector).WithPadding(-199)
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

func (s searchState) IntersectsFixed(page int, b primitives.BBox) bool {
	for i, v := range s.processed {
		if page != v.Page {
			continue
		}
		fixedBox := s.boxes[i].Translate(v.Vector).WithPadding(-199) // grow the box by the padding amount (minus a bit)
		// fmt.Printf("Box is %v, but padded is %v\n", box, fixedBox)
		if fixedBox.DoesIntersect(b) {
			// there is a conflict
			// fmt.Printf("There is an intersection between %v and %v\n", fixedBox, b)
			return true
		}
	}
	return false
}

// return the area of the bounding box that contains all positioned processed boxes
func (s searchState) ProcessedArea() float64 {
	pointsByPage := map[int][]primitives.Point{}
	for i, v := range s.processed {
		bx := s.boxes[i].Translate(v.Vector)
		pointsByPage[v.Page] = append(pointsByPage[v.Page], bx.BBox.Corners()...)
	}
	totalArea := 0.0
	for _, points := range pointsByPage {
		boundingBox := primitives.BBoxAroundPoints(points...)
		totalArea += boundingBox.Area()
	}
	return totalArea
}

type pageStat struct {
	page             int
	boxes            int
	bboxArea         float64
	points           []primitives.Point
	componentAreaSum float64
}

func (s searchState) PageAreaStats() map[int]pageStat {
	byPage := map[int]pageStat{}
	for i, v := range s.processed {
		page := v.Page
		if _, ok := byPage[page]; !ok {
			byPage[page] = pageStat{
				page:   page,
				points: []primitives.Point{},
			}
		}
		bx := s.boxes[i].Translate(v.Vector)
		stat := byPage[page]
		stat.points = append(stat.points, bx.BBox.Corners()...)
		stat.boxes += 1
		stat.componentAreaSum += bx.Area()
		byPage[page] = stat
	}
	for pageID, stats := range byPage {
		bbox := primitives.BBoxAroundPoints(stats.points...)
		stats.bboxArea = bbox.Area()
		byPage[pageID] = stats
	}
	return byPage
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

	for len(searchStates) > 0 {
		newStates := []*searchState{} // append to this slice, then swap this out with searchStates at the end
		fmt.Printf("have %d search states\n", len(searchStates))
		searchStates = consolidateSearchStates(searchStates)
		for _, state := range searchStates {
			// fmt.Printf("%v, %v\n", state.processed, state.positions)
			found := false
			if len(state.unprocessed) == 0 {
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
					if pos.IsUpperLeftCorner {
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
					newState = newState.RemoveSuperfluousPositions()
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
		}
		searchStates = newStates
	}
	if len(finalStates) == 0 {
		return PackingSolution{
			Pages:          0,
			Translations:   make(map[int]PagedVector, len(bxs)),
			DebugPositions: []PagedVector{},
		}
	}
	fmt.Printf("Got %d possible solutions\n", len(finalStates))
	solution := slices.MaxFunc(finalStates, func(a, b *searchState) int {
		if a.Pages() != b.Pages() {
			return cmp.Compare(a.Pages(), b.Pages())
		}
		// there should be as little unprocessable
		if len(a.unprocessables) != len(b.unprocessables) {
			return cmp.Compare(a.unprocessedArea(), b.unprocessedArea())
		}
		// sort by total area of processed
		return cmp.Compare(a.ProcessedArea(), b.ProcessedArea())
	})
	fmt.Printf("Best solution has %d unplaceable, %d placeable boxes\n", len(solution.unprocessables), len(solution.processed))
	fmt.Printf("Boxes are:\n")
	for i, v := range solution.processed {
		fmt.Printf("\t%v\n", boxes[i].Translate(v.Vector))
	}
	fmt.Printf("\nPositions are:\n")
	for _, pos := range solution.positions {
		fmt.Printf("\t%v\n", pos)
	}
	return PackingSolution{
		Pages:          solution.Pages(),
		Translations:   solution.processed,
		DebugPositions: solution.positions,
	}
}

func filterWithTooManyPages(states []*searchState) []*searchState {
	// remove solutions with too many pages
	// pages := map[int]int{}
	minpages := 10000
	maxpages := 1
	for _, state := range states {
		// pages[state.Pages()] += 1
		if state.Pages() < minpages {
			minpages = state.Pages()
		}
		if state.Pages() > maxpages {
			maxpages = state.Pages()
		}
	}
	fmt.Printf("\nmin: %d, max: %d\n", minpages, maxpages)
	newStates := []*searchState{}
	for _, state := range states {
		if state.Pages() <= minpages+1 {
			newStates = append(newStates, state)
		}
	}
	return newStates
}

func filterWithTooMuchUnusedSpace(states []*searchState) []*searchState {
	// consolidate states that leave 50%+ space empty
	newStates := []*searchState{}
	for _, state := range states {
		if len(state.processed) > 3 {
			pageStats := state.PageAreaStats()
			shouldRemove := false
			for _, stats := range pageStats {
				if stats.boxes > 3 {
					if stats.componentAreaSum/stats.bboxArea < 0.67 {
						shouldRemove = true
						break
					}
				}
			}
			if !shouldRemove { // TODO: dynamically determine the threshold
				newStates = append(newStates, state)
			}
		} else {
			newStates = append(newStates, state)
		}
	}
	return newStates
}

func filterWithSameFootprint(states []*searchState) []*searchState {
	// consolidate states that contain the same boxes and cover the same area
	stateMap := map[string]*searchState{}
	for _, state := range states {
		key := fmt.Sprintf("%s:%f", state.processedKey(), state.ProcessedArea())
		stateMap[key] = state
	}
	fmt.Printf("to %d ", len(stateMap))
	return maps.Values(stateMap)
}

func consolidateSearchStates(states []*searchState) []*searchState {
	if len(states) < 2 {
		return states
	}
	fmt.Printf("Consolidating from %d... ", len(states))
	states = filterWithTooManyPages(states)
	fmt.Printf("to %d ", len(states))

	states = filterWithTooMuchUnusedSpace(states)
	fmt.Printf("to %d ", len(states))

	states = filterWithSameFootprint(states)
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
