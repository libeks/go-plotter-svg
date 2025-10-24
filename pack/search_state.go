package pack

import (
	"cmp"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/libeks/go-plotter-svg/primitives"

	"github.com/kelindar/bitmap"
)

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

func (s searchState) RemoveSuperfluousPositions(container primitives.BBox) searchState {
	if len(s.positions) == 0 {
		return s
	}
	newPositions := []PagedVector{}
	for _, pos := range s.positions {
		pt := primitives.Origin.Add(pos.Vector)
		// skip if position falls outside of the container
		if !container.PointInside(pt) {
			continue
		}
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

// // return the area of the bounding box that contains all positioned processed boxes
// func (s searchState) ProcessedArea() float64 {
// 	pointsByPage := map[int][]primitives.Point{}
// 	for i, v := range s.processed {
// 		bx := s.boxes[i].Translate(v.Vector)
// 		pointsByPage[v.Page] = append(pointsByPage[v.Page], bx.BBox.Corners()...)
// 	}
// 	totalArea := 0.0
// 	for _, points := range pointsByPage {
// 		boundingBox := primitives.BBoxAroundPoints(points...)
// 		totalArea += boundingBox.Area()
// 	}
// 	return totalArea
// }

func stateComparatorFunc(a, b *searchState) int {
	// fewer pages is better
	if a.Pages() != b.Pages() {
		return cmp.Compare(a.Pages(), b.Pages())
	}
	// there should be as few unprocessable rectangles
	if a.unprocessables.Count() != b.unprocessables.Count() {
		return cmp.Compare(a.unprocessables.Count(), b.unprocessables.Count())
	}
	// sort by total area of processed
	return cmp.Compare(a.ProcessedAreaSum(), b.ProcessedAreaSum())
}

type pageStat struct {
	page             int
	nBoxes           int
	bboxArea         float64
	boxes            bitmap.Bitmap
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
				boxes:  bitmap.Bitmap{},
			}
		}
		bx := s.boxes[i].Translate(v.Vector)
		stat := byPage[page]
		stat.points = append(stat.points, bx.BBox.Corners()...)
		stat.nBoxes += 1
		stat.boxes.Set(uint32(i))
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

func (s searchState) BBoxAreaSum() float64 {
	total := 0.0
	stats := s.PageAreaStats()
	for _, stat := range stats {
		total += stat.bboxArea
	}
	return total
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
