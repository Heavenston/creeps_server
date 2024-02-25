package spatialevents

import (
	"lib.creeps.heav.fr/events"
	. "lib.creeps.heav.fr/geom"
	"lib.creeps.heav.fr/spatialmap"
)

type sub[T any] struct {
    filter AABB
    sendChan chan T
    handle *events.CancelHandle
}

func (sub sub[T]) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
    return nil
}

func (sub sub[T]) GetAABB() AABB {
    return sub.filter
}

// like events.EventProvider but filtered by position
type SpatialEventProvider[T spatialmap.Spatialized] struct {
    subs spatialmap.SpatialMap[sub[T]]
}

func NewSpatialEventProvider[T spatialmap.Spatialized]() *SpatialEventProvider[T] {
    this := new(SpatialEventProvider[T])
    return this
}

// if the filter has size 0 it matches every events
func (provider *SpatialEventProvider[T]) Subscribe(
    channel chan T,
    filter AABB,
) *events.CancelHandle {
    handle := new(events.CancelHandle)
    provider.SubscribeWithHandle(channel, filter, handle)
    return handle
}

// if the filter has size 0 it matches every events
func (provider *SpatialEventProvider[T]) SubscribeWithHandle(
    channel chan T,
    filter AABB,
    handle *events.CancelHandle,
) {
    if handle.IsCancelled() {
        return;
    }

    provider.subs.Add(sub[T]{
        sendChan: channel,
        handle: handle,
        filter: filter,
    })
}

func (provider *SpatialEventProvider[T]) Emit(event T) {
    aabb := event.GetAABB()

    provider.subs.RemoveAll(func(t sub[T]) bool {
        if t.handle.IsCancelled() {
            return true
        }

        if aabb.IsZero() || t.filter.IsZero() || t.filter.Intersects(aabb) {
            t.sendChan <- event
        }

        return false
    })
}
