package geom

import (
	"sync"

	mathutils "github.com/heavenston/creeps_server/creeps_lib/math_utils"
)

type AtomicPoint struct {
	lock  sync.RWMutex
	point Point
}

func (p *AtomicPoint) Load() Point {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.point
}

// Stores the given point and returns the previous value
func (p *AtomicPoint) Store(o Point) Point {
	p.lock.Lock()
	defer p.lock.Unlock()

	prev := p.point
	p.point = o
	return prev
}

// atomic modification of the position, returns the pervious values and new value
// respectively
func (p *AtomicPoint) Modify(cb func(Point) Point) (Point, Point) {
	p.lock.Lock()
	defer p.lock.Unlock()
	prev := p.point
	p.point = cb(p.point)
	return prev, p.point
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Point) ToAtomic() AtomicPoint {
	return AtomicPoint{
		point: p,
	}
}

func (a Point) Add(b Point) Point {
	return Point{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func (a Point) Plus(x int, y int) Point {
	return Point{
		X: a.X + x,
		Y: a.Y + y,
	}
}

func (a Point) Sub(b Point) Point {
	return Point{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func (a Point) Minus(x int, y int) Point {
	return Point{
		X: a.X - x,
		Y: a.Y - y,
	}
}

func (a Point) Times(v int) Point {
	return Point{
		X: a.X * v,
		Y: a.Y * v,
	}
}

func (a Point) Dist(b Point) int {
	return mathutils.AbsInt(a.X-b.X) + mathutils.AbsInt(a.Y-b.Y)
}
