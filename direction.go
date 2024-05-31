package main

var (
	Vertical   CardinalDirection = "vertical"
	Horizontal CardinalDirection = "horizontal"

	HomeToAway OrderDirection = "homeToAway"
	AwayToHome OrderDirection = "awayToHome"

	SameDirection        Connection = "sameDirection"
	AlternatingDirection Connection = "alternatingDirection"
	// ConnectedZigZag      Connection = "connectedZigZag" //unimplemented
)

type Direction struct {
	CardinalDirection                // which direction do all strokes point to (vertical, horizontal)
	StrokeDirection   OrderDirection // which direction should the first stroke be (home->away, away->home)
	OrderDirection                   // which direction should consecutive strokes go in (home->away, away->home)
	Connection                       // are strokes connected, do they alternate or all go in the same direction (sameDirection, oppositeDirection, connectedZigZag)

}

type CardinalDirection string

// type StrokeDirection string

type OrderDirection string

type Connection string
