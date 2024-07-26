package lines

import (
	"fmt"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/primitives"
)

type LineSegment struct {
	P1 primitives.Point
	P2 primitives.Point
}

func (l LineSegment) String() string {
	return fmt.Sprintf("LineSegment (%s) -> (%s)", l.P1, l.P2)
}

func (l LineSegment) Reverse() LineSegment {
	return LineSegment{l.P2, l.P1}
}

func (l LineSegment) XML(color, width string) xmlwriter.Elem {
	return xmlwriter.Elem{
		Name: "line", Attrs: []xmlwriter.Attr{
			{
				Name:  "x1",
				Value: fmt.Sprintf("%.1f", l.P1.X),
			},
			{
				Name:  "x2",
				Value: fmt.Sprintf("%.1f", l.P2.X),
			},
			{
				Name:  "y1",
				Value: fmt.Sprintf("%.1f", l.P1.Y),
			},
			{
				Name:  "y2",
				Value: fmt.Sprintf("%.1f", l.P2.Y),
			},
			{
				Name:  "stroke",
				Value: color,
			},
			{
				Name:  "fill",
				Value: "none",
			},
			{
				Name:  "stroke-width",
				Value: width,
			},
		},
	}
}

func (l LineSegment) GuideXML(color, width string) xmlwriter.Elem {
	return l.XML(color, width)
}

func (l LineSegment) Len() float64 {
	return l.P2.Subtract(l.P1).Len()
}

func (l LineSegment) Start() primitives.Point {
	return l.P1
}

func (l LineSegment) End() primitives.Point {
	return l.P2
}

func (l LineSegment) Line() Line {
	return Line{
		P: l.P1,
		V: l.P2.Subtract(l.P1),
	}
}

func (l LineSegment) IsEmpty() bool {
	return l.P1 == l.P2
}

func (l LineSegment) IntersectLineT(l2 Line) *float64 {
	l1 := l.Line()
	t := l1.IntersectT(l2)
	if t == nil {
		return nil
	}
	tt := *t
	if tt <= 1.0 && tt >= 0.0 {
		return t
	}
	return nil
}

func (l LineSegment) IntersectLineSegmentT(ls2 LineSegment) *float64 {
	l1 := l.Line()
	l2 := ls2.Line()
	t, u := l1.IntersectTU(l2)
	if t == nil || u == nil {
		return nil
	}
	tt := *t
	uu := *u
	if (tt <= 1.0 && tt >= 0.0) && (uu <= 1.0 && uu >= 0.0) {
		return t
	}
	return nil
}
