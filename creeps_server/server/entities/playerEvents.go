package entities

import (
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	. "github.com/heavenston/creeps_server/creeps_server/server"
)

type PlayerSpawnEvent struct {
	ServerEventBase
	Player *Player
}

func (event *PlayerSpawnEvent) GetAABB() AABB {
	// empty aabb = covers all map
	return AABB{}
}

type PlayerDespawnEvent struct {
	ServerEventBase
	Player *Player
}

func (event *PlayerDespawnEvent) GetAABB() AABB {
	// empty aabb = covers all map
	return AABB{}
}
