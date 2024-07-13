package truchet

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/samplers"
)

func TruchetTilesWithStrikeThrough(b box.Box, dataSource samplers.DataSource) []*Curve {
	val := dataSource.GetValue(b.Center())
	if val < 0.4 {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: West,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: East,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	} else if val > 0.6 {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: East,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: West,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	} else {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: LineOver,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: West,
						midpoint: 0.5,
					},
					{
						endpoint: East,
						midpoint: 0.5,
					},
				},
				CurveType: LineUnder,
			},
		}
	}
}

func TruchetTiles(b box.Box, dataSource samplers.DataSource) []*Curve {
	val := dataSource.GetValue(b.Center())
	if val < 0.5 {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: West,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: East,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	} else {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: East,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: West,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	}
}

type EndpointMidpoint struct {
	endpoint NWSE
	midpoint float64
}

func (e EndpointMidpoint) String() string {
	return fmt.Sprintf("%s %.1f", e.endpoint, e.midpoint)
}

type cellCoord struct {
	x int
	y int
}
