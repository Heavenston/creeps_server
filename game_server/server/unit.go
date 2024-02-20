package server

import (
	"fmt"

	"creeps.heav.fr/epita_api/model"
	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/spatialmap"
	"creeps.heav.fr/uid"
)

type UnitBusyError struct{}

func (e UnitBusyError) Error() string {
	return "unit is busy"
}

type UnsuportedActionError struct {
	Tried     ActionOpCode
	Supported []ActionOpCode
}

func (e UnsuportedActionError) Error() string {
	return fmt.Sprintf("action %s is not supported, supported: %v", e.Tried, e.Supported)
}

type NotEnoughResourcesError struct {
	Required model.Resources
	Available model.Resources
}

func (e NotEnoughResourcesError) Error() string {
	return fmt.Sprintf("not enough resources to perform action")
}

// every unit operation must be thread-safe atomic
// implemented in the server/units package (avoids circular depedency)
type IUnit interface {
	GetServer() *Server
	GetId() uid.Uid
	// returns an identifier of this kind of unit
	GetOpCode() string
	IsBusy() bool
	GetAlive() bool
	SetDead()
	// the id of the owner, note: can be the server by way of ServerUid
	GetOwner() uid.Uid
	GetPosition() Point
	SetPosition(newPos Point)
	GetAABB() AABB
	MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent]
	// atomically modifies the position of the unit
	ModifyPosition(cb func(Point) Point) (Point, Point)
	GetLastAction() *Action
	// can return UnitBusyError or UnsuportedActionError
	StartAction(action *Action) error
	GetUpgradeCosts() *model.CostResponse
	IsUpgraded() bool
	ObserveDistance() int
	GetInventory() model.Resources
	// atomically modifier the inventory
	ModifyInventory(func(model.Resources) model.Resources)
	SetInventory(newInv model.Resources)
	// Ran each tick after being registered by the server
	// only if GetAlive returns true
	Tick()
}
