package spatialevents_test

import (
	"testing"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	"github.com/heavenston/creeps_server/creeps_lib/events/spatialevents"
	"github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
	"github.com/stretchr/testify/assert"
)

type Event struct {
    extent spatialmap.Extent
}

func (self Event) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
    return nil
}

func (self Event) GetExtent() spatialmap.Extent {
    return self.extent
}

func extent(x int, y int, w int, h int, global bool) spatialmap.Extent {
    return spatialmap.Extent{
        IsGlobal: global,
        Aabb: geom.AABB{
            From: geom.Point{
                X: x,
                Y: y,
            },
            Size: geom.Point{
                X: w,
                Y: h,
            },
        },
    }
}

func event(extent spatialmap.Extent) Event {
    return Event { extent: extent }
}

func checkRecv[T comparable](t *testing.T, channel chan T, expect T) {
    t.Helper()
    select {
    case e, ok := (<- channel):
        if !ok {
            t.Errorf("Channel closed but expected event '%v'", expect)
        } else if e != expect {
            t.Errorf("Received '%v' but expected '%v'", e, expect)
        }
    default:
        t.Errorf("Received nothing but expected '%v'", expect)
    }
}

func checkEmpty[T comparable](t *testing.T, channel chan T) {
    t.Helper()
    select {
    case e, ok := (<- channel):
        if !ok {
            t.Errorf("Channel is closed but expected empty")
            return
        }
        t.Errorf("Received '%v' but expected empty", e)
    default:
    }
}

func checkClosed[T comparable](t *testing.T, channel chan T) {
    t.Helper()
    select {
    case e, ok := (<- channel):
        if !ok {
            break
        }
        t.Errorf("Shouldn't have received an event: %v", e)
    default:
        t.Errorf("Channel should've been closed")
    }
}

func TestSpatialBasic(t *testing.T) {
    provider := spatialevents.NewSpatialEventProvider[Event]()

    provider.Emit(event(extent(12, -12, 3, 3, false)))

    chan1 := make(chan Event, 10)
    cancel := provider.Subscribe(chan1, extent(10, -10, 20, 5, false))

    checkEmpty(t, chan1)

    provider.Emit(event(extent(-5, 80, 3, 3, false)))

    checkEmpty(t, chan1)

    ev1 := event(extent(5, -15, 10, 10, false))
    c1 := provider.Emit(ev1)
    assert.Equal(t, 1, c1)
    ev2 := event(extent(5, -15, 5, 5, false))
    c2 := provider.Emit(ev2)
    assert.Equal(t, 0, c2)
    ev3 := event(extent(5, -15, 15, 10, false))
    c3 := provider.Emit(ev3)
    assert.Equal(t, 1, c3)

    checkRecv(t, chan1, ev1)
    checkRecv(t, chan1, ev3)
    checkEmpty(t, chan1)

    cancel.Cancel()

    // channel only really closed after first emit
    // could be change but be carefull
    checkEmpty(t, chan1)
    
    provider.Emit(ev2)

    checkClosed(t, chan1)
}

func TestSpatialGlobalOnlyGlobal(t *testing.T) {
    provider := spatialevents.NewSpatialEventProvider[Event]()
    provider.Emit(event(extent(5, 15, 3, 3, false)))

    chan1 := make(chan Event, 10)
    cancel := provider.Subscribe(chan1, spatialmap.Extent{})

    checkEmpty(t, chan1)

    provider.Emit(event(extent(5, 80, 3, 3, false)))

    checkEmpty(t, chan1)

    ev1 := event(extent(-5, 80, 0, 0, true))
    provider.Emit(ev1)
    ev2 := event(extent(-5, 58, 1, 2, false))
    provider.Emit(ev2)
    ev3 := event(extent(-5, 58, 0, 0, true))
    provider.Emit(ev3)

    checkRecv(t, chan1, ev1)
    checkRecv(t, chan1, ev3)
    checkEmpty(t, chan1)

    cancel.Cancel()

    // channel only really closed after first emit
    // could be change but be carefull
    checkEmpty(t, chan1)
    
    provider.Emit(ev2)

    checkClosed(t, chan1)
}

func TestSpatialGlobalGlobal(t *testing.T) {
    provider := spatialevents.NewSpatialEventProvider[Event]()
    provider.Emit(event(extent(5, 15, 3, 3, false)))

    chan1 := make(chan Event, 10)
    cancel := provider.Subscribe(chan1, spatialmap.Extent{
        IsGlobal: true,
    })

    checkEmpty(t, chan1)

    ev0 := event(extent(5, 80, 3, 3, false))
    provider.Emit(ev0)

    checkRecv(t, chan1, ev0)

    ev1 := event(extent(-5, 80, 0, 0, true))
    provider.Emit(ev1)
    ev2 := event(extent(-5, 58, 1, 2, false))
    provider.Emit(ev2)
    ev3 := event(extent(-5, 58, 0, 0, true))
    provider.Emit(ev3)

    checkRecv(t, chan1, ev1)
    checkRecv(t, chan1, ev2)
    checkRecv(t, chan1, ev3)
    checkEmpty(t, chan1)

    cancel.Cancel()

    // channel only really closed after first emit
    // could be change but be carefull
    checkEmpty(t, chan1)
    
    provider.Emit(ev2)

    checkClosed(t, chan1)
}
