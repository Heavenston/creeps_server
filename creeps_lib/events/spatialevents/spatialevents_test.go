package spatialevents_test

import (
	"testing"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	"github.com/heavenston/creeps_server/creeps_lib/events/spatialevents"
	"github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
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

func event(x int, y int, w int, h int, global bool) Event {
    return Event {
        extent: spatialmap.Extent{
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
        },
    }
}

func checkRecv[T comparable](t *testing.T, channel chan T, expect T) {
    t.Helper()
    select {
    case e := (<- channel):
        if e != expect {
            t.Errorf("Wrong event received: %v", e)
        }
    default:
        t.Errorf("Should've received event '%v' but received nothing", expect)
    }
}

func checkEmpty[T comparable](t *testing.T, channel chan T) {
    t.Helper()
    select {
    case e, ok := (<- channel):
        if !ok {
            t.Errorf("Unexpected closed channel")
            return
        }
        t.Errorf("Unexpected event '%v'", e)
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

func TestPatialGlobalOnlyGlobal(t *testing.T) {
    provider := spatialevents.NewSpatialEventProvider[Event]()
    provider.Emit(event(5, 15, 3, 3, false))

    chan1 := make(chan Event, 10)
    cancel := provider.Subscribe(chan1, spatialmap.Extent{})

    checkEmpty(t, chan1)

    provider.Emit(event(5, 80, 3, 3, false))

    checkEmpty(t, chan1)

    ev1 := event(-5, 80, 0, 0, true)
    provider.Emit(ev1)
    ev2 := event(-5, 58, 1, 2, false)
    provider.Emit(ev2)
    ev3 := event(-5, 58, 0, 0, true)
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
