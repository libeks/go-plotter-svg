package truchet

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

func (c *Curve) XMLChunk(from endPointTuple) lines.PathChunk {
	if !c.HasEndpoint(from) {
		panic(fmt.Errorf("curve %s doesn't have endpoint %s", c, from))
	}
	to := c.GetOtherDirection(from)
	if to == nil {
		panic("No 'to' direction")
	}
	mTo := c.GetMidpoint(*to)
	mFrom := c.GetMidpoint(from)
	tFrom := *mFrom
	tTo := *mTo
	startPoint := c.Cell.AtEdge(from, tFrom)
	endPoint := c.Cell.AtEdge(*to, tTo)
	curveType := GetCurveType(from.NWSE, to.NWSE, tFrom, tTo)

	return c.Cell.Grid.CurveMapper.GetPathChunk(c, curveType, tFrom, tTo, startPoint, endPoint)
}

func (c *Curve) GetMidpoint(endpoint endPointTuple) *float64 {
	for _, pt := range c.endpoints {
		if pt.endpoint.endpoint == endpoint.endpoint {
			return &pt.midpoint
		}
	}
	return nil
}

func (c Curve) HasEndpoint(endpoint endPointTuple) bool {
	for _, pt := range c.endpoints {
		if pt.endpoint.endpoint == endpoint.endpoint {
			return true
		}
	}
	return false
}

// return the index of the other end of this curve, in corrdinates relative to this cell
func (c *Curve) GetOtherDirection(endpoint endPointTuple) *endPointTuple {
	var other *endPointTuple
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
