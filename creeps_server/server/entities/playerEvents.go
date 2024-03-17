package entities

import (
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
	. "github.com/heavenston/creeps_server/creeps_server/server"
)

type PlayerSpawnEvent struct {
	ServerEventBase
	Player *Player
}

func (event *PlayerSpawnEvent) GetExtent() spatialmap.Extent {
	return event.Player.GetExtent()
}

type PlayerDespawnEvent struct {
	ServerEventBase
	Player *Player
}

func (event *PlayerDespawnEvent) GetExtent() spatialmap.Extent {
	return event.Player.GetExtent()
}
