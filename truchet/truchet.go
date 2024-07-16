package truchet

import (
	"fmt"
)

func NewPair(a, b int) pair {
	return pair{a: a, b: b}
}

type pair struct {
	a int
	b int
}

func (p pair) Other(q int) int {
	if p.a == q {
		return p.b
	} else if p.b == q {
		return p.a
	}
	return -1
}

type tileSet struct {
	pairs []pair
}

func (t tileSet) Other(i int) int {
	for _, pair := range t.pairs {
		if other := pair.Other(i); other > 0 {
			return other
		}
	}
	return -1
}

// non-intersecting links for a 4-set, corresponds to Catalan number 2
var TruchetPairs = []tileSet{
	{
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 4),
		},
	},
	{
		pairs: []pair{
			NewPair(1, 4),
			NewPair(2, 3),
		},
	},
}

// all links for 4-set
var TruchetUnderPairs = []tileSet{
	{
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 4),
		},
	},
	{
		// cross
		pairs: []pair{
			NewPair(1, 3),
			NewPair(2, 4),
		},
	},
	{
		pairs: []pair{
			NewPair(1, 4),
			NewPair(2, 3),
		},
	},
}

// non-intersecting links for a 6-set, corresponds to Catalan number 3
var Truchet6Pairs = []tileSet{
	{
		// ()()()
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 4),
			NewPair(5, 6),
		},
	},
	{
		// ()(())
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 6),
			NewPair(4, 5),
		},
	},
	{
		// (())()
		pairs: []pair{
			NewPair(1, 4),
			NewPair(2, 3),
			NewPair(5, 6),
		},
	},
	{
		// ((()))
		pairs: []pair{
			NewPair(1, 6),
			NewPair(2, 5),
			NewPair(3, 4),
		},
	},
	{
		// (()())
		pairs: []pair{
			NewPair(1, 6),
			NewPair(2, 3),
			NewPair(4, 5),
		},
	},
}

type endPointTuple struct {
	endpoint int
	NWSE
}

func (e endPointTuple) isEmpty() bool {
	return e.endpoint == 0
}

type endPointPair struct {
	a endPointTuple
	b endPointTuple
}

func (p endPointPair) Has(q int) bool {
	if p.a.endpoint == q {
		return true
	} else if p.b.endpoint == q {
		return true
	}
	return false
}

func (p endPointPair) Other(q int) endPointTuple {
	if p.a.endpoint == q {
		return p.b
	} else if p.b.endpoint == q {
		return p.a
	}
	return endPointTuple{}
}

type edgePointMapping struct {
	pairs []endPointPair
}

func (e edgePointMapping) getHorizontal() []endPointPair {
	ret := []endPointPair{}
	for _, pair := range e.pairs {
		if pair.a.NWSE == North || pair.a.NWSE == South {
			ret = append(ret, pair)
		}
	}
	return ret
}

func (e edgePointMapping) getVertical() []endPointPair {
	ret := []endPointPair{}
	for _, pair := range e.pairs {
		if pair.a.NWSE == East || pair.a.NWSE == West {
			ret = append(ret, pair)
		}
	}
	return ret
}

func (e edgePointMapping) getDirection(i int) endPointTuple {
	for _, pair := range e.pairs {
		if pair.a.endpoint == i {
			return pair.a
		}
		if pair.b.endpoint == i {
			return pair.b
		}
	}
	return endPointTuple{}
}

func (e edgePointMapping) other(i int) endPointTuple {
	for _, pair := range e.pairs {
		other := pair.Other(i)
		if !other.isEmpty() {
			return other
		}
	}
	return endPointTuple{}
}

func (e edgePointMapping) endpointsFrom(direction NWSE) []endPointTuple {
	ret := []endPointTuple{}
	for _, pair := range e.pairs {
		if pair.a.NWSE == direction {
			ret = append(ret, pair.a)
		}
		if pair.b.NWSE == direction {
			ret = append(ret, pair.b)
		}
	}
	return ret
}

var EndpointMapping4 = edgePointMapping{
	[]endPointPair{
		{
			a: endPointTuple{
				endpoint: 1,
				NWSE:     North,
			},
			b: endPointTuple{
				endpoint: 3,
				NWSE:     South,
			},
		},
		{
			a: endPointTuple{
				endpoint: 2,
				NWSE:     East,
			},
			b: endPointTuple{
				endpoint: 4,
				NWSE:     West,
			},
		},
	},
}

var EndpointMapping6Side = edgePointMapping{
	[]endPointPair{
		{
			a: endPointTuple{
				endpoint: 1,
				NWSE:     North,
			},
			b: endPointTuple{
				endpoint: 4,
				NWSE:     South,
			},
		},
		{
			a: endPointTuple{
				endpoint: 2,
				NWSE:     East,
			},
			b: endPointTuple{
				endpoint: 6,
				NWSE:     West,
			},
		},
		{
			a: endPointTuple{
				endpoint: 3,
				NWSE:     East,
			},
			b: endPointTuple{
				endpoint: 5,
				NWSE:     West,
			},
		},
	},
}

type EndpointMidpoint struct {
	endpoint endPointTuple
	midpoint float64
}

func (e EndpointMidpoint) String() string {
	return fmt.Sprintf("%s %.1f", e.endpoint, e.midpoint)
}

type cellCoord struct {
	x int
	y int
}
