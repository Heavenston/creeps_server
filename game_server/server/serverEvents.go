package server

import (
	mathutils "creeps.heav.fr/math_utils"
	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/spatialmap"
)

// utility struct embedded into all server events to auto-implement the functions
type serverEventBase struct {
	
}

func (event *serverEventBase) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return nil
}

// see spatialmap.Spatialized
type IServerEvent interface {
	MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent]
	GetAABB() AABB
}

// emitted by the server on the call of RegisterUnit
type UnitSpawnEvent struct {
	serverEventBase
	Unit IUnit
	AABB AABB
}

func (event *UnitSpawnEvent) GetAABB() AABB {
	return event.AABB
}

// emitted by the server on the call of RemoveUnit
type UnitDespawnEvent struct {
	serverEventBase
	Unit IUnit
	AABB AABB
}

func (event *UnitDespawnEvent) GetAABB() AABB {
	return event.AABB
}

// emitted by the units when setPosition is called
type UnitMovedEvent struct {
	serverEventBase
	Unit IUnit
	From Point
	To   Point
}

func (event *UnitMovedEvent) GetAABB() AABB {
	// FIXME: This could take much more space than what's really required
	//        (aabb for From and one for To) but as the current design doesn't
	//        support finer control I'll stick with this as most movements
	//        are across adjacent tiles anyways

	min := Point {
		X: mathutils.Min(event.From.X, event.To.X),
		Y: mathutils.Min(event.From.Y, event.To.Y),
	}
	max := Point {
		X: mathutils.Max(event.From.X, event.To.X),
		Y: mathutils.Max(event.From.Y, event.To.Y),
	}

	return AABB{
		From: min,
		Size: max.Sub(min),
	}
}

// emitted by the server on the call of RegisterPlayer
type PlayerSpawnEvent struct {
	serverEventBase
	Player *Player
}

func (event *PlayerSpawnEvent) GetAABB() AABB {
	// empty aabb = covers all map
	return AABB{}
}

// emitted by the server on the call of RemovePlayer
type PlayerDespawnEvent struct {
	serverEventBase
	Player *Player
}

func (event *PlayerDespawnEvent) GetAABB() AABB {
	// empty aabb = covers all map
	return AABB{}
}
