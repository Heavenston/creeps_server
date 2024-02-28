package spatialevents

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
	"github.com/rs/zerolog/log"
)

type sub[T any] struct {
    filter AABB
    sendChan chan T
    handle *events.CancelHandle

    file string
    line int
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

    _, file, line, _ := runtime.Caller(1)
    if strings.Contains(file, "spatialevents") {
        _, file, line, _ = runtime.Caller(2)
    }

    provider.subs.Add(sub[T]{
        sendChan: channel,
        handle: handle,
        filter: filter,
        file: file,
        line: line,
    })
}

func (provider *SpatialEventProvider[T]) Emit(event T) {
    aabb := event.GetAABB()

    provider.subs.RemoveAll(func(t sub[T]) bool {
        return t.handle.IsCancelled()
    })

    for _, sub := range provider.subs.GetAllIntersects(aabb) {
        select {
        case sub.sendChan <- event:
        default:
            log.Warn().
                Str("event_type", reflect.TypeOf(event).String()).
                Str("sub_file", sub.file).
                Int("sub_line", sub.line).
                Msg("Could not send event")
        }
    }
}
