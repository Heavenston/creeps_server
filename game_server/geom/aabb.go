package geom

// An AxisAlignedBoundingBox, size must be position or else functions will return invalid results
// (but size ~can~ be 0 -> contains no point and not contained in anything)
type AABB struct {
	// minimum of both axis
	From Point
	Size Point
}

// returns true if the size is zero
func (aabb AABB) IsZero() bool {
	return aabb.Size.X == 0 && aabb.Size.Y == 0
}

func (aabb AABB) GetPosition() Point {
	return aabb.From
}

// gets the maximum value (excluded) of both axis
func (aabb AABB) Upto() Point {
	return aabb.From.Add(aabb.Size)
}

// returns true if the given point is inside the aabb
func (aabb AABB) Contains(p Point) bool {
	return p.X >= aabb.From.X && p.Y >= aabb.From.Y &&
		p.X < aabb.Upto().X && p.Y < aabb.Upto().Y
}

// retuns true if
func (aabb AABB) Intersects(other AABB) bool {
	return aabb.From.X < other.Upto().X && aabb.From.Y < other.Upto().Y &&
		other.From.X < aabb.Upto().X && other.From.Y < aabb.Upto().Y
}
