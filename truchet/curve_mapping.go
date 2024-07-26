package truchet

import (
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
)

var (
	BlockyCurveMapper = CurveMapper{
		straightLineMapper,
		quarterBezierLineMapper,
		loopbackLineMapper,
	}

	MapCircularCurve = CurveMapper{
		straightLineMapper,
		quarterCircleBezLineMapper,
		loopbackLineMapper,
	}

	MapCurlyCurve = CurveMapper{
		straightLineMapper,
		quarterBezLineMapper,
		loopbackLineMapper,
	}
)

func straightLineMapper(c *Curve, curveType CurveType, tFrom, tTo float64, startPoint, endPoint primitives.Point) lines.PathChunk {
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
}

func loopbackLineMapper(c *Curve, curveType CurveType, tFrom, tTo float64, startPoint, endPoint primitives.Point) lines.PathChunk {
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

func quarterBezierLineMapper(c *Curve, curveType CurveType, tFrom, tTo float64, startPoint, endPoint primitives.Point) lines.PathChunk {
	crossesDiagonal := c.GetClockIntersectDiagonal(curveType, tFrom, tTo)

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
}

func quarterCircleBezLineMapper(c *Curve, curveType CurveType, tFrom, tTo float64, startPoint, endPoint primitives.Point) lines.PathChunk {
	crossesDiagonal := c.GetClockIntersectDiagonal(curveType, tFrom, tTo)
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
			// TODO: redo this logic to not have to have such odd conditionals inside matching switch cases
			if curveType == ClockNE {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tFrom, maths.Interpolate(0, tTo, 0.55)),
					P2:    c.Cell.At(maths.Interpolate(tFrom, 1, 0.55), tTo),
					End:   endPoint,
				}
			}
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(tFrom, maths.Interpolate(1, tTo, 0.55)),
				P2:    c.Cell.At(maths.Interpolate(0, tFrom, 0.55), tTo),
				End:   endPoint,
			}
		}

	case CClockEN, CClockWS:
		if crossesDiagonal {
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(tFrom, tFrom),
				P2:    c.Cell.At(tTo, tTo),
				End:   endPoint,
			}
		} else {
			if curveType == CClockEN {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(maths.Interpolate(1, tTo, 0.55), tFrom),
					P2:    c.Cell.At(tTo, maths.Interpolate(0, tFrom, 0.55)),
					End:   endPoint,
				}
			}
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(maths.Interpolate(0, tTo, 0.55), tFrom),
				P2:    c.Cell.At(tTo, maths.Interpolate(1, tFrom, 0.55)),
				End:   endPoint,
			}
		}

	case ClockWN, ClockES:
		if crossesDiagonal {
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(1-tFrom, tFrom),
				P2:    c.Cell.At(tTo, 1-tTo),
				End:   endPoint,
			}
		} else {
			if curveType == ClockWN {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(maths.Interpolate(0, tTo, 0.55), tFrom),
					P2:    c.Cell.At(tTo, maths.Interpolate(tFrom, 0, 0.55)),
					End:   endPoint,
				}
			}
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(maths.Interpolate(1, tTo, 0.55), tFrom),
				P2:    c.Cell.At(tTo, maths.Interpolate(1, tFrom, 0.55)),
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
			if curveType == CClockNW {
				return lines.CubicBezierChunk{
					Start: startPoint,
					P1:    c.Cell.At(tFrom, maths.Interpolate(0, tTo, 0.55)),
					P2:    c.Cell.At(maths.Interpolate(tFrom, 0, 0.55), tTo),
					End:   endPoint,
				}
			}
			return lines.CubicBezierChunk{
				Start: startPoint,
				P1:    c.Cell.At(tFrom, maths.Interpolate(1, tTo, 0.55)),
				P2:    c.Cell.At(maths.Interpolate(1, tFrom, 0.55), tTo),
				End:   endPoint,
			}
		}
	default:
		panic("Unexpected case")
	}
}

func quarterBezLineMapper(c *Curve, curveType CurveType, tFrom, tTo float64, startPoint, endPoint primitives.Point) lines.PathChunk {
	switch curveType {
	case ClockNE, ClockSW:
		return lines.CubicBezierChunk{
			Start: startPoint,
			P1:    c.Cell.At(tFrom, tFrom),
			P2:    c.Cell.At(tTo, tTo),
			End:   endPoint,
		}
	case CClockEN, CClockWS:
		return lines.CubicBezierChunk{
			Start: startPoint,
			P1:    c.Cell.At(tFrom, tFrom),
			P2:    c.Cell.At(tTo, tTo),
			End:   endPoint,
		}

	case ClockWN, ClockES:
		return lines.CubicBezierChunk{
			Start: startPoint,
			P1:    c.Cell.At(1.0-tFrom, tFrom),
			P2:    c.Cell.At(tTo, 1.0-tTo),
			End:   endPoint,
		}
	case CClockNW, CClockSE:
		return lines.CubicBezierChunk{
			Start: startPoint,
			P1:    c.Cell.At(tFrom, 1.0-tFrom),
			P2:    c.Cell.At(1.0-tTo, tTo),
			End:   endPoint,
		}
	default:
		panic("Unexpected case")
	}
}

type curveComponentMapper func(*Curve, CurveType, float64, float64, primitives.Point, primitives.Point) lines.PathChunk

type CurveMapper struct {
	straightLineMapper curveComponentMapper
	quarterMapper      curveComponentMapper
	loopbackLineMapper curveComponentMapper
}

func (m CurveMapper) GetPathChunk(c *Curve, curveType CurveType, tFrom, tTo float64, startPoint, endPoint primitives.Point) lines.PathChunk {
	switch curveType.MetaType() {
	case Straight:
		return m.straightLineMapper(c, curveType, tFrom, tTo, startPoint, endPoint)
	case QuarterCircle:
		return m.quarterMapper(c, curveType, tFrom, tTo, startPoint, endPoint)
	case Loopback:
		return m.loopbackLineMapper(c, curveType, tFrom, tTo, startPoint, endPoint)
	default:
		panic("Unexpected case")
	}
}
