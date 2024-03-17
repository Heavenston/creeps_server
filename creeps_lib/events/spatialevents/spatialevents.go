package spatialevents

import (
	"runtime"
	"strings"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
	"github.com/rs/zerolog/log"
)

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
	extent spatialmap.Extent,
) *events.CancelHandle {
	handle := new(events.CancelHandle)
	provider.SubscribeWithHandle(channel, extent, handle)
	return handle
}

// if the filter has size 0 it matches every events
func (provider *SpatialEventProvider[T]) SubscribeWithHandle(
	channel chan T,
	extent spatialmap.Extent,
	handle *events.CancelHandle,
) {
	if handle.IsCancelled() {
		return
	}

	_, file, line, _ := runtime.Caller(1)
	if strings.Contains(file, "spatialevents") {
		_, file, line, _ = runtime.Caller(2)
	}

	provider.subs.Add(sub[T]{
		sendChan: channel,
		handle:   handle,
		file:     file,
		line:     line,
	})
}

func (provider *SpatialEventProvider[T]) Emit(event T) {
	provider.subs.RemoveAll(func(t sub[T]) bool {
		if t.handle.IsCancelled() {
			close(t.sendChan)
			return true
		}
		return false
	})

	for _, sub := range provider.subs.GetAllCollides(event.GetExtent()) {
		select {
		case sub.sendChan <- event:
		default:
			log.Warn().
				Type("event_type", event).
				Str("sub_file", sub.file).
				Int("sub_line", sub.line).
				Msg("Could not send event")
		}
	}
}
