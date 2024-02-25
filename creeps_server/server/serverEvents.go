package server

import (
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/events"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	mathutils "github.com/heavenston/creeps_server/creeps_lib/math_utils"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
)

// utility struct embedded into all server events to auto-implement the functions
type ServerEventBase struct {
}

func (event *ServerEventBase) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return nil
}

// see spatialmap.Spatialized
type IServerEvent interface {
	MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent]
	GetAABB() AABB
}

type UnitSpawnEvent struct {
	ServerEventBase
	Unit IUnit
	AABB AABB
}

func (event *UnitSpawnEvent) GetAABB() AABB {
	return event.AABB
}

type UnitDespawnEvent struct {
	ServerEventBase
	Unit IUnit
	AABB AABB
}

func (event *UnitDespawnEvent) GetAABB() AABB {
	return event.AABB
}

type UnitMovedEvent struct {
	ServerEventBase
	Unit IUnit
	From Point
	To   Point
}

func (event *UnitMovedEvent) GetAABB() AABB {
	// FIXME: This could take much more space than what's really required
	//        (aabb for From and one for To) but as the current design doesn't
	//        support finer control I'll stick with this as most movements
	//        are across adjacent tiles anyways

	min := Point{
		X: mathutils.Min(event.From.X, event.To.X),
		Y: mathutils.Min(event.From.Y, event.To.Y),
	}
	max := Point{
		X: mathutils.Max(event.From.X, event.To.X),
		Y: mathutils.Max(event.From.Y, event.To.Y),
	}

	return AABB{
		From: min,
		Size: max.Sub(min),
	}
}

type UnitStartedActionEvent struct {
	ServerEventBase
	Unit   IUnit
	Pos    Point
	Action *Action
}

func (event *UnitStartedActionEvent) GetAABB() AABB {
	return AABB{
		From: event.Pos,
		Size: Point{X: 1, Y: 1},
	}
}

type UnitFinishedActionEvent struct {
	ServerEventBase
	Unit   IUnit
	Pos    Point
	Action *Action
	Report model.IReport
}

func (event *UnitFinishedActionEvent) GetAABB() AABB {
	return AABB{
		From: event.Pos,
		Size: Point{X: 1, Y: 1},
	}
}
