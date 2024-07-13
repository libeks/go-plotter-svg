package truchet

type Winding int

const (
	Clockwise Winding = iota
	CounterClockwise
	Straight
	LoopBack
	Undefined
)

func (w Winding) String() string {
	return []string{"Clockwise", "CounterClockwise", "Straight", "LoopBack", "Undefined"}[w]
}
