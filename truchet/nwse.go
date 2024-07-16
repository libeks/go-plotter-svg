package truchet

import "fmt"

type NWSE int

const (
	North NWSE = iota
	West
	South
	East

	UnknownNWSE
)

func (d NWSE) String() string {
	return []string{"North", "West", "South", "East", "UnknownNWSE"}[d]
}

func (d NWSE) Opposite() NWSE {
	switch d {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	default:
		panic(fmt.Errorf("direction %s doesn't have an opposite", d))
	}
}

func (d NWSE) CurveMetaType(next NWSE) CurveMetaType {
	if d == next {
		return Loopback
	}
	switch d {
	case North:
		switch next {
		case West:
			return QuarterCircle
		case South:
			return Straight
		case East:
			return QuarterCircle
		}
	case East:
		switch next {
		case North:
			return QuarterCircle
		case West:
			return Straight
		case South:
			return QuarterCircle
		}
	case South:
		switch next {
		case East:
			return QuarterCircle
		case North:
			return Straight
		case West:
			return QuarterCircle
		}
	case West:
		switch next {
		case South:
			return QuarterCircle
		case East:
			return Straight
		case North:
			return QuarterCircle
		}
	}
	return UnknownCurveMetaType
}
