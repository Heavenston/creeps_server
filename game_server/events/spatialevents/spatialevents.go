package spatialevents

import (
	"sync"

	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/spatialmap"
)

type sub[T any] struct {
    filter AABB
    sendChan chan T
    handle *events.CancelHandle
}

// like events.EventProvider but filtered by position
// FIXME: (fixme,,, or not) it's upsetting its the same code as EventProvider
//        but with a filter, technically with this implementation it could
//        be generic by using a function as filter but the idea was that
//        it could be optimized for 2d positions (so this is to revisit later)
type SpatialEventProvider[T spatialmap.Spatialized] struct {
    mutex sync.Mutex
    subs []sub[T]
}

func NewSpatialEventProvider[T spatialmap.Spatialized]() *SpatialEventProvider[T] {
    this := new(SpatialEventProvider[T])
    return this
}

func (provider *SpatialEventProvider[T]) Subscribe(
    channel chan T,
    filter AABB,
) *events.CancelHandle {
    handle := new(events.CancelHandle)
    provider.SubscribeWithHandle(channel, filter, handle)
    return handle
}

func (provider *SpatialEventProvider[T]) SubscribeWithHandle(
    channel chan T,
    filter AABB,
    handle *events.CancelHandle,
) {
    if handle.IsCancelled() {
        return;
    }

    provider.mutex.Lock()
    defer provider.mutex.Unlock()

    provider.subs = append(provider.subs, sub[T]{
        sendChan: channel,
        handle: handle,
        filter: filter,
    })
}

func (provider *SpatialEventProvider[T]) Emit(event T) {
    aabb := event.GetAABB()
    
    var i int = 0
    for i < len(provider.subs) {
        sub := &provider.subs[i]
        if sub.handle.IsCancelled() {
            close(sub.sendChan)
            copy(provider.subs[i:], provider.subs[i+1:])
            provider.subs = provider.subs[:len(provider.subs)-1]
            continue
        }

        if aabb.IsZero() || aabb.Intersects(sub.filter) {
            sub.sendChan <- event
        }

        i++
    }
}
