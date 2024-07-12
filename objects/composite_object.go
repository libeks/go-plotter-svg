package objects

import (
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

func NewCompositeWithWithout(with []Object, without []Object) CompositeObject {
	return CompositeObject{}.With(with...).Without(without...)
}

func NewComposite() CompositeObject {
	return CompositeObject{}
}

type CompositeObject struct {
	positive []Object
	negative []Object
}

func (o CompositeObject) With(obj ...Object) CompositeObject {
	return CompositeObject{
		positive: append(o.positive, obj...),
		negative: o.negative,
	}
}

func (o CompositeObject) Without(obj ...Object) CompositeObject {
	return CompositeObject{
		positive: o.positive,
		negative: append(o.negative, obj...),
	}
}

func (o CompositeObject) Inside(p primitives.Point) bool {
	inside := false
	for _, pos := range o.positive {
		if pos.Inside(p) {
			inside = true
			break
		}
	}
	if !inside {
		return false
	}
	for _, neg := range o.negative {
		if neg.Inside(p) {
			return false
		}
	}
	return true
}

func (o CompositeObject) IntersectTs(line lines.Line) []float64 {
	ts := []float64{}
	for _, obj := range append(o.positive, o.negative...) {
		ts = append(ts, obj.IntersectTs(line)...)
	}
	return ts
}

func (o CompositeObject) IntersectCircleTs(circle Circle) []float64 {
	ts := []float64{}
	for _, obj := range append(o.positive, o.negative...) {
		ts = append(ts, obj.IntersectCircleTs(circle)...)
	}
	return ts
}
