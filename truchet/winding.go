package truchet

type Winding int

const (
	Clockwise Winding = iota
	CounterClockwise
	Straight
	Undefined
)

func (w Winding) String() string {
	return []string{"Clockwise", "CounterClockwise", "Straight", "Undefined"}[w]
}
