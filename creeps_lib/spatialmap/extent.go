package spatialmap

import (
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
)

type Extent struct {
	Aabb     AABB
	// If set to true, collides with all extents regardless of its aabb
	IsGlobal bool
}

func (self Extent) Collides(other Extent) bool {
	if other.IsGlobal || self.IsGlobal {
		return true
	}
	return self.Aabb.Intersects(other.Aabb)
}

func (self Extent) Contains(point Point) bool {
	if self.IsGlobal {
		return true
	}
	return self.Aabb.Contains(point)
}
