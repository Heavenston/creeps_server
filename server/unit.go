package server

import (
	"fmt"

	"creeps.heav.fr/api/model"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/uid"
)

type UnitBusyError struct{}

func (e UnitBusyError) Error() string {
	return "unit is busy"
}

type UnsuportedActionError struct {
	Tried     ActionOpCode
}

func (e UnsuportedActionError) Error() string {
	return fmt.Sprintf("action %s is not supported", e.Tried)
}

// every unit operation must be thread-safe atomic
type IUnit interface {
	GetServer() *Server
	GetId() uid.Uid
	GetAlive() bool
	SetAlive(new bool)
	// the id of the owner, note: can be the server by way of ServerUid
	GetOwner() uid.Uid
	GetPosition() Point
	SetPosition(newPos Point)
	// atomically modifies the position of the unit
	ModifyPosition(cb func(Point) Point) (Point, Point)
	GetLastAction() *Action
	StartAction(action *Action) error
	GetUpgradeCosts() *model.CostResponse
	// Ran each tick after being registered by the server
	// only if GetAlive returns true
	Tick()
}
