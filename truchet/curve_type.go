package truchet

type CurveMetaType int
type CurveType int

// is this enough? should I have a method to convert frmo some abstraction of the curve to its representation?
// so there is
// from (point, direction and t-value)
// to (point, direction, t-value)
// - Straight
// VerticalDown
// VerticalUp
// HorizontalLeft
// HorizontalRight
// - QuarterCircle
// ClockNE
// ClockES
// ClockSW
// ClockWN
// CClockNW
// CClockWS
// CClockSE
// CClockEN
// - Loopbacks
// LoopbackEUp
// LoopbackEDown
// LoopbackWUp
// LoopbackWDown
// LoopbackNLeft
// LoopbackNRight
// LoopbackSLeft
// LoopbackSRight

const (
	Straight CurveMetaType = iota
	QuarterCircle
	Loopback

	UnknownCurveMetaType
)

func (c CurveMetaType) String() string {
	return []string{"Straight", "QuarterCircle", "Loopback", "UnknownCurveMetaType"}[c]
}

func GetCurveMetaType(from, to NWSE) CurveMetaType {
	if from == to {
		return Loopback
	}
	switch from {
	case North:
		switch to {
		case West:
			return QuarterCircle
		case South:
			return Straight
		case East:
			return QuarterCircle
		}
	case East:
		switch to {
		case North:
			return QuarterCircle
		case West:
			return Straight
		case South:
			return QuarterCircle
		}
	case South:
		switch to {
		case East:
			return QuarterCircle
		case North:
			return Straight
		case West:
			return QuarterCircle
		}
	case West:
		switch to {
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

const (
	// Straight
	VerticalDown CurveType = iota
	VerticalUp
	HorizontalLeft
	HorizontalRight
	// QuarterCircle
	ClockNE
	ClockES
	ClockSW
	ClockWN
	CClockNW
	CClockWS
	CClockSE
	CClockEN
	// Loopbacks
	LoopbackEUp
	LoopbackEDown
	LoopbackWUp
	LoopbackWDown
	LoopbackNLeft
	LoopbackNRight
	LoopbackSLeft
	LoopbackSRight

	UnknownCurveType
)

func (c CurveType) String() string {
	return []string{"VerticalDown", "VerticalUp", "HorizontalLeft", "HorizontalRight", "ClockNE", "ClockES", "ClockSW", "ClockWN", "CClockNW", "CClockWS", "CClockSE", "CClockEN", "LoopbackEUp", "LoopbackEDown", "LoopbackWUp", "LoopbackWDown", "LoopbackNLeft", "LoopbackNRight", "LoopbackSLeft", "LoopbackSRight", "UnknownCurveType"}[c]
}

func (c CurveType) MetaType() CurveMetaType {
	switch c {
	case VerticalDown, VerticalUp, HorizontalLeft, HorizontalRight:
		return Straight
	case ClockNE, ClockES, ClockSW, ClockWN, CClockNW, CClockWS, CClockSE, CClockEN:
		return QuarterCircle
	case LoopbackEUp, LoopbackEDown, LoopbackWUp, LoopbackWDown, LoopbackNLeft, LoopbackNRight, LoopbackSLeft, LoopbackSRight:
		return Loopback
	default:
		return UnknownCurveMetaType
	}
}

func GetCurveType(from, to NWSE, fromT, toT float64) CurveType {
	if from == to {
		fromBeforeTo := fromT < toT
		switch from {
		case East:
			if fromBeforeTo {
				return LoopbackEDown
			} else {
				return LoopbackEUp
			}
		case West:
			if fromBeforeTo {
				return LoopbackWDown
			} else {
				return LoopbackWUp
			}
		case North:
			if fromBeforeTo {
				return LoopbackNRight
			} else {
				return LoopbackNLeft
			}
		case South:
			if fromBeforeTo {
				return LoopbackSRight
			} else {
				return LoopbackSLeft
			}
		}
		return UnknownCurveType
	}
	switch from {
	case North:
		switch to {
		case West:
			return CClockNW
		case South:
			return VerticalDown
		case East:
			return ClockNE
		}
	case East:
		switch to {
		case North:
			return CClockEN
		case West:
			return HorizontalLeft
		case South:
			return ClockES
		}
	case South:
		switch to {
		case East:
			return CClockSE
		case North:
			return VerticalUp
		case West:
			return ClockSW
		}
	case West:
		switch to {
		case South:
			return CClockWS
		case East:
			return HorizontalRight
		case North:
			return ClockWN
		}
	}
	return UnknownCurveType
}
