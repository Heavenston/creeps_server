package spatialmap

import "creeps.heav.fr/geom"

type Positioned interface {
	comparable
	GetPosition() geom.Point
}

type SpatialMap[T Positioned] struct {
	objects []T
}

func (m *SpatialMap[T]) Add(p T) {
	for _, o := range m.objects {
		if o == p {
			panic("cannot have duplicate objects")
		}
	}
	m.objects = append(m.objects, p)
}

func (m *SpatialMap[T]) RemoveFirst(predicate func(T) bool) *T {
	for i, o := range m.objects {
		if predicate(o) {
			m.objects[i] = m.objects[len(m.objects)-1]
			m.objects = m.objects[:len(m.objects)-1]
			return &o
		}
	}
	return nil
}

func (m *SpatialMap[T]) Find(predicate func(T) bool) *T {
	for _, obj := range m.objects {
		if predicate(obj) {
			return &obj
		}
	}
	return nil
}

func (m *SpatialMap[T]) GetAt(point geom.Point) *T {
	return m.Find(func(t T) bool {
		return t.GetPosition() == point
	})
}

func (m *SpatialMap[T]) GetIn(from geom.Point, upto geom.Point) (*T, bool) {
	for _, obj := range m.objects {
		if obj.GetPosition().IsWithing(from, upto) {
			return &obj, true
		}
	}
	return nil, false
}

func (m *SpatialMap[T]) Iter() func() (bool, int, *T) {
	var i int = 0
	return (func() (bool, int, *T) {
		if i == len(m.objects) {
			return false, 0, nil
		}
		v := m.objects[i]
		i += 1
		return true, i - 1, &v
	})
}
