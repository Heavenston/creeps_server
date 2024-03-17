package spatialmap_test

import (
	"slices"
	"testing"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	"github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
)

type Obj struct {
	aabb   geom.AABB
	events *events.EventProvider[spatialmap.ObjectMovedEvent]
}

func (self Obj) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return self.events
}

func (self Obj) GetAABB() geom.AABB {
	return self.aabb
}

func aabb(x int, y int, w int, h int) geom.AABB {
	return geom.AABB{
		From: geom.Point{
			X: x,
			Y: y,
		},
		Size: geom.Point{
			X: w,
			Y: h,
		},
	}
}

func obj(x int, y int, w int, h int) *Obj {
	return &Obj{
		events: nil,
		aabb:   aabb(x, y, w, h),
	}
}

func TestSpatialMapGetAt(t *testing.T) {
	sm := spatialmap.NewSpatialMap[*Obj]()
	defer sm.Close()

	obj1 := obj(0, 0, 10, 10)
	sm.Add(obj1)
	obj2 := obj(5, 5, 10, 5)
	sm.Add(obj2)
	obj3 := obj(50, -10, 10, 5)
	sm.Add(obj3)

	found := sm.GetAt(geom.Point{X: 5, Y: 5})
	if found == nil {
		t.Fatalf("No obj found")
	}
	f := *found
	if f != obj1 && f != obj2 {
		t.Fatalf("Wrong obj found")
	}

	found = sm.GetAt(geom.Point{X: 0, Y: 0})
	if found == nil {
		t.Fatalf("No obj found")
	}
	f = *found
	if f != obj1 {
		t.Fatalf("Wrong obj found")
	}

	found = sm.GetAt(geom.Point{X: 51, Y: -6})
	if found == nil {
		t.Fatalf("No obj found")
	}
	f = *found
	if f != obj3 {
		t.Fatalf("Wrong obj found")
	}
}

// Like slices.Equal but slices order is not significant
func slicesEquiv[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !slices.Contains(a, b[i]) {
			return false
		}
		if !slices.Contains(b, a[i]) {
			return false
		}
	}

	return true
}

func checkSlicesEquiv[T comparable](t *testing.T, got []T, expected []T) {
	t.Helper()

	if len(got) != len(expected) {
		t.Errorf(
			"Slices length don't match, got %d, expected %d",
			len(got),
			len(expected),
		)
	}

	for i := 0; i < len(got) || i < len(expected); i++ {
		if i < len(expected) && !slices.Contains(got, expected[i]) {
			t.Errorf("'%v' is absent", expected[i])
		}
		if i < len(got) && !slices.Contains(expected, got[i]) {
			t.Errorf("'%v' is extra", got[i])
		}
	}
}

func TestSpatialMapGetIntersects(t *testing.T) {
	sm := spatialmap.NewSpatialMap[Obj]()
	defer sm.Close()

	obj1 := *obj(0, 0, 10, 10)
	sm.Add(obj1)
	obj2 := *obj(5, 5, 10, 5)
	sm.Add(obj2)
	obj3 := *obj(50, -10, 10, 5)
	sm.Add(obj3)

	checkSlicesEquiv(t, sm.GetAllIntersects(aabb(0, 0, 0, 0)), []Obj{})
	checkSlicesEquiv(t, sm.GetAllIntersects(aabb(52, 99, 0, 0)), []Obj{})
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(-50, -50, 100, 100)),
		[]Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(-50, -50, 200, 200)),
		[]Obj{obj1, obj2, obj3},
	)
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(0, 0, 50, 50)),
		[]Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(5, 5, 5, 5)),
		[]Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(0, 0, 3, 3)),
		[]Obj{obj1},
	)
}

func TestSpatialMapGlobalObjects(t *testing.T) {
	sm := spatialmap.NewSpatialMap[Obj]()
	defer sm.Close()

	obj1 := *obj(0, 0, 0, 0)
	sm.Add(obj1)
	obj2 := *obj(-2, 5, 0, 0)
	sm.Add(obj2)
	obj3 := *obj(-5, 5, 1, 2)
	sm.Add(obj3)
	obj4 := *obj(1, 5, 3, 2)
	sm.Add(obj4)

	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(0, 0, 0, 0)), []Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(-8, 0, 8, 10)), []Obj{obj1, obj2, obj3},
	)
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(-10, -10, 30, 30)), []Obj{obj1, obj2, obj3, obj4},
	)
	checkSlicesEquiv(t,
		sm.GetAllIntersects(aabb(0, 5, 10, 10)), []Obj{obj1, obj2, obj4},
	)
}
