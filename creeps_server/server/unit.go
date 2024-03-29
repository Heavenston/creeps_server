package server

import (
	"fmt"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

type UnitBusyError struct{}

func (e UnitBusyError) Error() string {
	return "unit is busy"
}

type UnsuportedActionError struct {
	Tried     model.ActionOpCode
	Supported []model.ActionOpCode
}

func (e UnsuportedActionError) Error() string {
	return fmt.Sprintf("action %s is not supported, supported: %v", e.Tried, e.Supported)
}

// every unit operation must be thread-safe atomic
// implemented in the server/units package (avoids circular depedency)
type IUnit interface {
	IEntity
	// returns an identifier of this kind of unit
	GetOpCode() string
	IsBusy() bool
	// the id of the owner, note: can be the server by way of ServerUid
	GetOwner() uid.Uid
	GetPosition() Point
	SetPosition(newPos Point)
	// atomically modifies the position of the unit
	ModifyPosition(cb func(Point) Point) (Point, Point)
	GetLastAction() *Action
	// can return UnitBusyError or UnsuportedActionError
	// unFinished can be nil
	StartAction(action *Action, onFinished func()) error
	GetUpgradeCosts() *model.CostResponse
	IsUpgraded() bool
	ObserveDistance() int
	GetInventory() model.Resources
	// atomically modifier the inventory
	ModifyInventory(func(model.Resources) model.Resources)
	SetInventory(newInv model.Resources)
}
