package curve

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/lines"
)

type Curve struct {
	*Cell
	endpoints []EndpointMidpoint
	visited   bool
}

func (c *Curve) String() string {
	return fmt.Sprintf("Curve at %s with endpoints %v", c.Cell, c.endpoints)
}

func (c *Curve) GetClockIntersectDiagonal(curveType CurveType, from, to float64) bool {
	// does the quarter circle cross the diagonal? if so, we need to do a cubic bezier curve that is
	// bounded by that diagonal. If not, we can do a simple quadratic bezier
	// WN or NW
	// +------+
	// |   | /|
	// |   |/ |
	// |   /  |
	// |--/   |
	// | /    |
	// |/     |
	// +------+

	// WS or SW
	// +------+
	// |\     |
	// |-\    |
	// |  \   |
	// |   \  |
	// |   |\ |
	// |   | \|
	// +------+

	// EN or NE
	// +------+  .0
	// |\|    |
	// | \    |  |
	// |  \   |  |
	// |   \--|  \/
	// |    \ |
	// |     \|
	// +------+  1.0
	// .0 -->  1.0
	//

	// tFrom, tFrom -> tTo, tTo

	switch curveType {
	case ClockWN, CClockNW:

		// NW or WN, order doesn't matter
		if from+to <= 1.0 {
			// do quadratic
			return false
		} else {
			// do cubic
			return true
		}
	case ClockNE, CClockWS:

		// NE, if from = 0.2 and to = 0.6 -> 1.0 + 0.2 - 0.6 = 0.6, doesn't qualify as cubic
		// WS, if from = 0.2 and to = 0.6 -> 1.0 - 0.2 + 0.6 = 1.4, qualifies as cubic

		// WS or SW
		if 1.0-from+to <= 1.0 {
			return false
		} else {
			return true
		}

	case ClockSW, CClockEN:

		// SW, if from = 0.6 and to = 0.2 -> 1.0 - 0.6 + 0.2 = 0.6, doesn't qualify as cubic
		// EN, if from = 0.6 and to = 0.2 -> 1.0 - 0.2 + 0.6 = 1.4, qualifies as cubic

		// NE or EN
		if 1.0-to+from <= 1.0 {
			return false
		} else {
			return true
		}

	case ClockES, CClockSE:
		// SE or ES
		// equivalent to 1.0 - tFrom + 1.0 - tTo < 1.0 <=> 2.0 - (tFrom + tTo) < 1.0 <=> tFrom + tTo > 1.0
		// again, order doesn't matter, it's symmetric
		if from+to >= 1.0 {
			return false
		} else {
			return true
		}
	default:
		panic("Unexpected curve type")
	}
}

func (c *Curve) XMLChunk(curveMapper CurveMapper, from connectionEnd) lines.PathChunk {
	if !c.HasEndpoint(from) {
		panic(fmt.Errorf("curve %s doesn't have endpoint %s", c, from))
	}
	to := c.GetOtherEnd(from)
	if to == nil {
		panic("No 'to' direction")
	}
	mTo := c.GetTValue(*to)
	mFrom := c.GetTValue(from)
	tFrom := *mFrom
	tTo := *mTo
	startPoint := c.Cell.AtEdge(from.NWSE, tFrom)
	if !c.Cell.PointInside(startPoint) {
		panic("startpoint not inside cell")
	}
	endPoint := c.Cell.AtEdge(to.NWSE, tTo)
	if !c.Cell.PointInside(endPoint) {
		panic("startpoint not inside cell")
	}
	curveType := GetCurveType(from.NWSE, to.NWSE, tFrom, tTo)

	return curveMapper.GetPathChunk(c, curveType, tFrom, tTo, startPoint, endPoint)
}

func (c *Curve) GetTValue(endpoint connectionEnd) *float64 {
	for _, pt := range c.endpoints {
		if pt.endpoint.endpoint == endpoint.endpoint {
			return &pt.tValue
		}
	}
	return nil
}

func (c Curve) HasEndpoint(endpoint connectionEnd) bool {
	for _, pt := range c.endpoints {
		if pt.endpoint.endpoint == endpoint.endpoint {
			return true
		}
	}
	return false
}

// return the index of the other end of this curve within the same cell
func (c *Curve) GetOtherEnd(endpoint connectionEnd) *connectionEnd {
	var other *connectionEnd
	found := false
	for _, pt := range c.endpoints {
		if pt.endpoint.endpoint == endpoint.endpoint {
			found = true
		} else {
			other = &pt.endpoint
		}
	}
	if found {
		return other
	}
	return nil
}

type EndpointMidpoint struct {
	endpoint connectionEnd
	tValue   float64
}

func (e EndpointMidpoint) String() string {
	return fmt.Sprintf("%s %.1f", e.endpoint, e.tValue)

}
