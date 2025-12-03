package truchet 


// connection represents the indexed point of a curve fragment within a cell
type connectionEnd struct {
	endpoint int
	NWSE
}

func (e connectionEnd) isEmpty() bool {
	return e.endpoint == 0
}

// connectionPair represents a connection across two adjoining cells
type connectionPair struct {
	a connectionEnd
	b connectionEnd
}

// Has returns true if this pair contains the indexed endpoint
func (p connectionPair) Has(q int) bool {
	if p.a.endpoint == q {
		return true
	} else if p.b.endpoint == q {
		return true
	}
	return false
}

func (p connectionPair) Other(q int) connectionEnd {
	if p.a.endpoint == q {
		return p.b
	} else if p.b.endpoint == q {
		return p.a
	}
	return connectionEnd{}
}

func (p connectionPair) bothEnds() []connectionEnd {
	return []connectionEnd{p.a, p.b}
}

// edgePointMapping defines how the same endpoint maps from one cell onto the next
type edgePointMapping struct {
	pairs []connectionPair
}

func (e edgePointMapping) getHorizontal() []connectionPair {
	ret := []connectionPair{}
	for _, pair := range e.pairs {
		if pair.a.NWSE == North || pair.a.NWSE == South {
			ret = append(ret, pair)
		}
	}
	return ret
}

func (e edgePointMapping) getVertical() []connectionPair {
	ret := []connectionPair{}
	for _, pair := range e.pairs {
		if pair.a.NWSE == East || pair.a.NWSE == West {
			ret = append(ret, pair)
		}
	}
	return ret
}

func (e edgePointMapping) getDirection(i int) connectionEnd {
	for _, pair := range e.pairs {
		if pair.a.endpoint == i {
			return pair.a
		}
		if pair.b.endpoint == i {
			return pair.b
		}
	}
	return connectionEnd{}
}

func (e edgePointMapping) other(i int) connectionEnd {
	for _, pair := range e.pairs {
		other := pair.Other(i)
		if !other.isEmpty() {
			return other
		}
	}
	return connectionEnd{}
}

func (e edgePointMapping) endpointsFrom(direction NWSE) []connectionEnd {
	ret := []connectionEnd{}
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