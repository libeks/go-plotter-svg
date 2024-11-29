package samplers

import "github.com/libeks/go-plotter-svg/primitives"

type Booleaner interface {
	GetBool(p primitives.Point) bool
}

type Not struct {
	Booleaner
}

func (n Not) GetBool(p primitives.Point) bool {
	return !n.Booleaner.GetBool(p)
}

type And struct {
	P1 Booleaner
	P2 Booleaner
}

func (a And) GetBool(p primitives.Point) bool {
	if !a.P1.GetBool(p) {
		return false
	}
	return a.P2.GetBool(p)
}

type Or struct {
	P1 Booleaner
	P2 Booleaner
}

func (b Or) GetBool(p primitives.Point) bool {
	if b.P1.GetBool(p) {
		return true
	}
	return b.P2.GetBool(p)
}

type Xor struct {
	P1 Booleaner
	P2 Booleaner
}

func (b Xor) GetBool(p primitives.Point) bool {
	p1 := b.P1.GetBool(p)
	p2 := b.P2.GetBool(p)
	if p1 && p2 {
		return false
	}
	if p1 || p2 {
		return true
	}
	return false
}

type ConcentricCircleBoolean struct {
	Center primitives.Point
	// true inside 0-radii[0], alternating afterwards
	Radii []float64
}

func (b ConcentricCircleBoolean) GetBool(p primitives.Point) bool {
	rad := CircleRadius{Center: b.Center}.GetValue(p)
	inside := true
	for _, radComp := range b.Radii {
		if rad < radComp {
			return inside
		}
		inside = !inside
	}
	return inside
}
