package entities

import (
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
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
