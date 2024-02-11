package geom

import (
	mathutils "creeps.heav.fr/math_utils"
)

type Point struct {
    X, Y int
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
    return mathutils.AbsInt(a.X - b.X) + mathutils.AbsInt(a.Y - b.Y)
}
