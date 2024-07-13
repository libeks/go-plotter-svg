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

func (d NWSE) Winding(next NWSE) Winding {
	switch d {
	case North:
		switch next {
		case West:
			return CounterClockwise
		case South:
			return Straight
		case East:
			return Clockwise
		}
	case East:
		switch next {
		case North:
			return CounterClockwise
		case West:
			return Straight
		case South:
			return Clockwise
		}
	case South:
		switch next {
		case East:
			return CounterClockwise
		case North:
			return Straight
		case West:
			return Clockwise
		}
	case West:
		switch next {
		case South:
			return CounterClockwise
		case East:
			return Straight
		case North:
			return Clockwise
		}
	}
	return Undefined
}

func (d NWSE) String() string {
	return []string{"North", "West", "South", "East", "UnknownNWSE"}[d]
}
