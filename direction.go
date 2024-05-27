package main

var (
	RightToLeft StrokeDirection = "rightToLeft"
	LeftToRight StrokeDirection = "leftToRight"

	TopToBottom OrderDirection = "topToBottom"
	BottomToTop OrderDirection = "bottomToTop"

	SameDirection     Connection = "sameDirection"
	OppositeDirection Connection = "oppositeDirection"
	ZigZag            Connection = "zigZag"
)

type Direction struct {
	StrokeDirection
	OrderDirection
	Connection
}

type StrokeDirection string

type OrderDirection string

type Connection string
