package truchet

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/lines"
)

type Curve struct {
	*Cell
	// endpoints NWSE
	// midpoints [2]float64 // in the range of [0;1], in increasing order of N,W,S,E
	endpoints []EndpointMidpoint
	CurveType
	visited bool
}

func (c *Curve) String() string {
	return fmt.Sprintf("Curve at %s with endpoints %v", c.Cell, c.endpoints)
}

func (c *Curve) XMLChunk(from NWSE) lines.PathChunk {
	if !c.HasEndpoint(from) {
		panic(fmt.Errorf("curve %s doesn't have endpoint %s", c, from))
	}
	to := c.GetOtherDirection(from)
	if to == nil {
		panic("No 'to' direction")
	}
	mTo := c.GetMidpoint(*to)
	mFrom := c.GetMidpoint(from)
	startPoint := c.Cell.At(from, *mFrom)
	endPoint := c.Cell.At(*to, *mTo)
	radius := c.Cell.Box.Width() / 2
	winding := from.Winding(*to)
	switch winding {
	case Straight:
		if c.CurveType == LineOver {
			fmt.Printf("doing line over %s\n", c)
			return lines.LineChunk{
				End: endPoint,
			}
		} else if c.CurveType == LineUnder {
			fmt.Printf("doing line under %s\n", c)
			return lines.LineGapChunk{
				Start:        startPoint,
				GapSizeRatio: 0.5,
				End:          endPoint,
			}
		} else {
			fmt.Printf("curve type %s\n", c.CurveType)
		}
	case Clockwise:
		return lines.CircleArcChunk{
			Radius:      radius,
			IsClockwise: false, // Truchet circle arcs swing the other direction from winding
			IsLong:      false,
			End:         endPoint,
		}
	case CounterClockwise:
		return lines.CircleArcChunk{
			Radius:      radius,
			IsClockwise: true, // Truchet circle arcs swing the other direction from winding
			IsLong:      false,
			End:         endPoint,
		}
	case Undefined:
		fmt.Printf("winding is undefined: %s\n", winding)
		return lines.LineChunk{
			End: endPoint,
		}
	default:
		fmt.Printf("winding is %s\n", winding)
		return lines.LineChunk{
			End: endPoint,
		}
	}
	fmt.Printf("not even default: winding is %s\n", winding)

	return lines.LineChunk{
		End: endPoint,
	}
}

func (c *Curve) GetMidpoint(endpoint NWSE) *float64 {
	for _, pt := range c.endpoints {
		if pt.endpoint == endpoint {
			return &pt.midpoint
		}
	}
	return nil
}

func (c Curve) HasEndpoint(endpoint NWSE) bool {
	for _, pt := range c.endpoints {
		if pt.endpoint == endpoint {
			return true
		}
	}
	return false
}

func (c *Curve) GetOtherDirection(endpoint NWSE) *NWSE {
	var other *NWSE
	found := false
	for _, pt := range c.endpoints {
		if pt.endpoint == endpoint {
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
