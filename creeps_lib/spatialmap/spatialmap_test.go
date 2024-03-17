package spatialmap_test

import (
	"slices"
	"testing"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	"github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
)

type Obj struct {
	extent   spatialmap.Extent
	events *events.EventProvider[spatialmap.ObjectMovedEvent]
}

func (self Obj) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return self.events
}

func (self Obj) GetExtent() spatialmap.Extent {
	return self.extent
}

func extent(x int, y int, w int, h int, global bool) spatialmap.Extent {
	aabb := geom.AABB{
		From: geom.Point{
			X: x,
			Y: y,
		},
		Size: geom.Point{
			X: w,
			Y: h,
		},
	}
	return spatialmap.Extent{
		IsGlobal: global,
		Aabb: aabb,
	}
}

func obj(x int, y int, w int, h int, global bool) *Obj {
	return &Obj{
		events: nil,
		extent:   extent(x, y, w, h, global),
	}
}

func TestSpatialMapGetAt(t *testing.T) {
	sm := spatialmap.NewSpatialMap[*Obj]()
	defer sm.Close()

	obj1 := obj(0, 0, 10, 10, false)
	sm.Add(obj1)
	obj2 := obj(5, 5, 10, 5, false)
	sm.Add(obj2)
	obj3 := obj(50, -10, 10, 5, false)
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

func TestSpatialMapGetCollides(t *testing.T) {
	sm := spatialmap.NewSpatialMap[Obj]()
	defer sm.Close()

	obj1 := *obj(0, 0, 10, 10, false)
	sm.Add(obj1)
	obj2 := *obj(5, 5, 10, 5, false)
	sm.Add(obj2)
	obj3 := *obj(50, -10, 10, 5, false)
	sm.Add(obj3)

	checkSlicesEquiv(t, sm.GetAllCollides(extent(0, 0, 0, 0, false)), []Obj{})
	checkSlicesEquiv(t, sm.GetAllCollides(extent(52, 99, 0, 0, false)), []Obj{})
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(-50, -50, 100, 100, false)),
		[]Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(-50, -50, 200, 200, false)),
		[]Obj{obj1, obj2, obj3},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(0, 0, 50, 50, false)),
		[]Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(5, 5, 5, 5, false)),
		[]Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(0, 0, 3, 3, false)),
		[]Obj{obj1},
	)
}

func TestSpatialMapGetCollides2(t *testing.T) {
	sm := spatialmap.NewSpatialMap[Obj]()
	defer sm.Close()

    obj1 := *obj(5, -15, 10, 10, false)
	sm.Add(obj1)
    obj2 := *obj(5, -15, 5, 5, false)
	sm.Add(obj2)
    obj3 := *obj(5, -15, 15, 10, false)
	sm.Add(obj3)

	checkSlicesEquiv(t,
	    sm.GetAllCollides(extent(10, -10, 20, 5, false)),
		[]Obj{obj1, obj3},
	)
}

func TestSpatialMapGlobalObjects(t *testing.T) {
	sm := spatialmap.NewSpatialMap[Obj]()
	defer sm.Close()

	obj1 := *obj(0, 0, 0, 0, true)
	sm.Add(obj1)
	obj2 := *obj(-2, 5, 0, 0, true)
	sm.Add(obj2)
	obj3 := *obj(-5, 5, 1, 2, false)
	sm.Add(obj3)
	obj4 := *obj(1, 5, 3, 2, false)
	sm.Add(obj4)

	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(0, 0, 0, 0, false)),
		[]Obj{obj1, obj2},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(-8, 0, 8, 10, false)),
		[]Obj{obj1, obj2, obj3},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(-10, -10, 30, 30, false)),
		[]Obj{obj1, obj2, obj3, obj4},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(0, 5, 10, 10, false)),
		[]Obj{obj1, obj2, obj4},
	)
}

func TestSpatialMapGetGlobal(t *testing.T) {
	sm := spatialmap.NewSpatialMap[Obj]()
	defer sm.Close()

	obj1 := *obj(0, 0, 0, 0, false)
	sm.Add(obj1)
	obj2 := *obj(-2, 5, 0, 0, false)
	sm.Add(obj2)
	obj3 := *obj(-5, 5, 1, 2, false)
	sm.Add(obj3)
	obj4 := *obj(1, 5, 3, 2, false)
	sm.Add(obj4)

	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(0, 0, 0, 0, true)),
		[]Obj{obj1, obj2, obj3, obj4},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(-8, 0, 8, 10, true)),
		[]Obj{obj1, obj2, obj3, obj4},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(-10, -10, 30, 30, true)),
		[]Obj{obj1, obj2, obj3, obj4},
	)
	checkSlicesEquiv(t,
		sm.GetAllCollides(extent(0, 5, 10, 10, true)),
		[]Obj{obj1, obj2, obj3, obj4},
	)
}
