package truchet

type CurveType int

const (
	StraightLine CurveType = iota
	CircleSegment
	Bezier
	LineOver
	LineUnder
	LoopBack
	Unknown
)

func (c CurveType) String() string {
	return []string{"StraightLine", "CircleSegment", "Bezier", "LineOver", "LineUnder", "LoopBack", "Unknown"}[c]
}
