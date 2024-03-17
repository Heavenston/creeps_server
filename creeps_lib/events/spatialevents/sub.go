package spatialevents

import (
	"github.com/heavenston/creeps_server/creeps_lib/events"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
)

type sub[T any] struct {
	extent   spatialmap.Extent
	sendChan chan T
	handle   *events.CancelHandle

	file string
	line int
}

func (sub sub[T]) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return nil
}

func (sub sub[T]) GetExtent() spatialmap.Extent {
	return sub.extent
}
