package geom

import (
	mathutils "creeps.heav.fr/math_utils"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (a Point) Add(b Point) Point {
	return Point{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func (a Point) Sub(b Point) Point {
	return Point{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func (a Point) Dist(b Point) int {
	return mathutils.AbsInt(a.X-b.X) + mathutils.AbsInt(a.Y-b.Y)
}

// Returns true if the point is within the rectangle delimited by the given points
// with from being included and upto expluded
func (p Point) IsWithing(from Point, upto Point) bool {
	min_x := mathutils.MinInt(from.X, upto.X)
	min_y := mathutils.MinInt(from.Y, upto.Y)
	max_x := mathutils.MaxInt(from.X, upto.X)
	max_y := mathutils.MaxInt(from.Y, upto.Y)

	return p.X >= min_x && p.Y >= min_y && p.X < max_x && p.Y < max_y
}
