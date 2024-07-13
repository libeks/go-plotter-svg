package truchet

type CurveType int

const (
	StraightLine CurveType = iota
	CircleSegment
	LineOver
	LineUnder
	LoopBack
	Unknown
)

func (c CurveType) String() string {
	return []string{"StraightLine", "CircleSegment", "LineOver", "LineUnder"}[c]
}
