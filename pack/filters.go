package pack

import (
	"fmt"

	// "golang.org/x/exp/maps"

	"github.com/libeks/go-plotter-svg/primitives"
)

func filterWithTooManyPages(states []*searchState) []*searchState {
	// remove solutions with too many pages
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
	newStates := []*searchState{}
	for _, state := range states {
		if state.Pages() <= minpages+1 {
			newStates = append(newStates, state)
		}
	}
	return newStates
}

func filterWithIncreasingBBoxesOnPages(states []*searchState) []*searchState {
	// consolidate states have later pages cover more area than earlier
	newStates := []*searchState{}
	for _, state := range states {
		pageStats := state.PageAreaStats()
		shouldRemove := false
		area := 10000000000.0 // some arbitrarily big number to begin with
		for pageID := range len(pageStats) {
			if pageStats[pageID].bboxArea > area {
				shouldRemove = true
				break
			}
			area = pageStats[pageID].bboxArea
		}
		if !shouldRemove {
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
				if stats.nBoxes > 3 {
					if stats.componentAreaSum/stats.bboxArea < 0.70 {
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

func filterFirstPageTooSmall(states []*searchState, container primitives.BBox) []*searchState {
	// consolidate states that leave 50%+ space empty
	newStates := []*searchState{}
	for _, state := range states {
		pageStats := state.PageAreaStats()
		shouldRemove := false
		// fmt.Printf("page stats %v, %f\n", pageStats, pageStats[0].bboxArea/container.Area())
		if len(pageStats) > 1 && pageStats[0].bboxArea < container.Area()/2 {
			shouldRemove = true
		}
		if !shouldRemove { // TODO: dynamically determine the threshold
			newStates = append(newStates, state)
		}
	}
	return newStates
}

func filterTwoPagesAddToLessThanContainer(states []*searchState, container primitives.BBox) []*searchState {
	// consolidate states where two or more pages have a bbox that could be placed on a single page
	newStates := []*searchState{}
	for _, state := range states {
		pageStats := state.PageAreaStats()
		shouldRemove := false
		// fmt.Printf("page stats %v, %f\n", pageStats, pageStats[0].bboxArea/container.Area())
		if len(pageStats) > 1 {
			total := 0.0
			for pageID, stat := range pageStats {
				if pageID > 0 && total+stat.bboxArea/container.Area() < 1.0 {
					shouldRemove = true
				}
				total += stat.bboxArea / container.Area()
			}
		}
		if !shouldRemove { // TODO: dynamically determine the threshold
			newStates = append(newStates, state)
		}
	}
	return newStates
}

func filterWithSameFootprintButMoreArea(states []*searchState) []*searchState {
	type areaStruct struct {
		area        float64
		searchState *searchState
	}
	// consolidate states that contain the same boxes and cover the same area
	stateMap := map[string]areaStruct{}
	for _, state := range states {
		stats := state.PageAreaStats()
		// TODO: see if the same set of objects on this page take up more area
		key := ""
		totalArea := 0.0
		for page, stat := range stats {
			st, err := stat.boxes.MarshalJSON()
			if err != nil {
				panic("Couldn't marshal")
			}
			key = fmt.Sprintf("%s:%d-%s", key, page, st)
			totalArea += stat.componentAreaSum
		}
		if val, ok := stateMap[key]; !ok {
			stateMap[key] = areaStruct{
				area:        totalArea,
				searchState: state,
			}
		} else {
			if val.area > totalArea {
				stateMap[key] = areaStruct{
					area:        totalArea,
					searchState: state,
				}
			}
		}
	}
	states = []*searchState{}
	for _, val := range stateMap {
		states = append(states, val.searchState)
	}
	// fmt.Printf("to %d ", len(stateMap))
	// return maps.Values(stateMap)
	return states
}

func filterSameBoxesWithMoreArea(states []*searchState) []*searchState {
	type areaStruct struct {
		area        float64
		searchState *searchState
	}
	// consolidate states that contain the same boxes to only keep the one that covers the least area
	stateMap := map[string]areaStruct{}
	for _, state := range states {
		area := state.BBoxAreaSum()
		key := state.processedKey()
		if alternate, ok := stateMap[key]; ok {
			if alternate.area > area {
				stateMap[key] = areaStruct{
					area:        area,
					searchState: state,
				}
			}
		} else {
			stateMap[key] = areaStruct{
				area:        area,
				searchState: state,
			}
		}
	}
	states = make([]*searchState, 0, len(stateMap))
	for _, value := range stateMap {
		states = append(states, value.searchState)
	}
	return states
}

func consolidateSearchStates(states []*searchState, container primitives.BBox) []*searchState {
	if len(states) < 2 {
		return states
	}
	fmt.Printf("Consolidating\n  from %d... \n", len(states))

	states = filterWithSameFootprintButMoreArea(states)
	fmt.Printf("  to %d (same footprint, more area)\n", len(states))

	states = filterWithTooManyPages(states)
	fmt.Printf("  to %d (too many pages)\n", len(states))

	// states = filterFirstPageTooSmall(states, container)
	// fmt.Printf("  to %d (first page too small)\n", len(states))

	states = filterTwoPagesAddToLessThanContainer(states, container)
	fmt.Printf("  to %d (multiple pages can be combined into one)\n", len(states))

	states = filterWithIncreasingBBoxesOnPages(states)
	fmt.Printf("  to %d (pages must be decreasing)\n", len(states))

	// TODO: reinstate this back once it does a statistical approach
	// states = filterWithTooMuchUnusedSpace(states)
	// fmt.Printf("  to %d (too much unused space)\n", len(states))

	// this doesn't seem to help at all
	// states = filterWithSameFootprint(states)
	// fmt.Printf("to %d ", len(states))

	// states = filterSameBoxesWithMoreArea(states)
	// fmt.Printf("  to %d (remove states with more area)\n", len(states))

	// states = maps.Values(stateMap)
	// // remove states that are identical
	// stateMap = map[string]searchState{}
	// for _, state := range states {
	// 	// fmt.Printf("key %s\n", state.Key())
	// 	stateMap[state.Key()] = state
	// }
	// fmt.Printf("to %d ", len(stateMap))
	// states = maps.Values(stateMap)

	// fmt.Printf("\n")
	return states
}
