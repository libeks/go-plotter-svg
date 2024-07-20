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
		if from+to < 1.0 {
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
		if 1.0-from+to < 1.0 {
			return false
		} else {
			return true
		}

	case ClockSW, CClockEN:

		// SW, if from = 0.6 and to = 0.2 -> 1.0 - 0.6 + 0.2 = 0.6, doesn't qualify as cubic
		// EN, if from = 0.6 and to = 0.2 -> 1.0 - 0.2 + 0.6 = 1.4, qualifies as cubic

		// NE or EN
		if 1.0-to+from < 1.0 {
			return false
		} else {
			return true
		}

	case ClockES, CClockSE:
		// SE or ES
		// equivalent to 1.0 - tFrom + 1.0 - tTo < 1.0 <=> 2.0 - (tFrom + tTo) < 1.0 <=> tFrom + tTo > 1.0
		// again, order doesn't matter, it's symmetric
		if from+to > 1.0 {
			return false
		} else {
			return true
		}
	default:
		panic("Unexpected curve type")
	}
}

// func (c *Curve) GetCurveType()

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

	switch curveType.MetaType() {
	case Straight:
		switch curveType {
		case HorizontalLeft, HorizontalRight:
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(0.5, tFrom),
				P2:    c.Cell.At(0.5, tTo),
				End:   endPoint,
			}
		case VerticalUp, VerticalDown:
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(tFrom, 0.5),
				P2:    c.Cell.At(tTo, 0.5),
				End:   endPoint,
			}
		default:
			panic("Unexpected case")
		}
	case QuarterCircle:
		crossesDiagonal := c.GetClockIntersectDiagonal(curveType, tFrom, tTo)
		fmt.Printf("type %s, crossesDiagonal %v, from %.1f to %.1f\n", curveType, crossesDiagonal, tFrom, tTo)

		switch curveType {
		// if diagonal does from top left to bottom right
		case ClockNE, ClockSW:
			if crossesDiagonal {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tFrom, tFrom),
					P2:    c.Cell.At(tTo, tTo),
					End:   endPoint,
				}
			} else {
				return lines.QuadraticBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tFrom, tTo), // this depends on the direction, only applies if start is N or S
					End:   endPoint,
				}
			}

		case CClockEN, CClockWS:
			fmt.Printf("type %s, crossesDiagonal %v, from %.1f to %.1f\n", curveType, crossesDiagonal, tFrom, tTo)
			if crossesDiagonal {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tFrom, tFrom),
					P2:    c.Cell.At(tTo, tTo),
					End:   endPoint,
				}
			} else {
				return lines.QuadraticBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tTo, tFrom), // this depends on the direction, only applies if start is E or W
					End:   endPoint,
				}
			}

		case ClockWN, ClockES:
			if crossesDiagonal {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(1.0-tFrom, tFrom),
					P2:    c.Cell.At(tTo, 1.0-tTo),
					End:   endPoint,
				}
			} else {
				return lines.QuadraticBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tTo, tFrom),
					End:   endPoint,
				}
			}
		case CClockNW, CClockSE:
			if crossesDiagonal {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tFrom, 1.0-tFrom),
					P2:    c.Cell.At(1.0-tTo, tTo),
					End:   endPoint,
				}
			} else {
				return lines.QuadraticBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tFrom, tTo),
					End:   endPoint,
				}
			}
		default:
			panic("Unexpected case")
		}
	case Loopback:
		switch curveType {
		case LoopbackWUp, LoopbackWDown:
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(0.3, tFrom),
				P2:    c.Cell.At(0.3, tTo),
				End:   endPoint,
			}
		case LoopbackEUp, LoopbackEDown:
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(0.7, tFrom),
				P2:    c.Cell.At(0.7, tTo),
				End:   endPoint,
			}
		case LoopbackNLeft, LoopbackNRight:
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(tFrom, 0.3),
				P2:    c.Cell.At(tTo, 0.3),
				End:   endPoint,
			}
		case LoopbackSLeft, LoopbackSRight:
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(tFrom, 0.7),
				P2:    c.Cell.At(tTo, 0.7),
				End:   endPoint,
			}
		default:
			panic("Unexpected case")
		}
	}

	return lines.LineChunk{
		End: endPoint,
	}
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
