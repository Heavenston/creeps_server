package units

import (
	"sync"
	"sync/atomic"

	"creeps.heav.fr/api/model"
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
)

type extendedUnit interface {
	IUnit
	getUnit() *unit
}

// See server.go's IUnit interface to explain its functions
type unit struct {
	// read-only (no lock)
	server *Server
	// read-only (no lock)
	id            uid.Uid
	alive         atomic.Bool
	position      AtomicPoint
	lastAction    atomic.Pointer[Action]
	upgraded      atomic.Bool
	inventoryLock sync.RWMutex
	inventory     model.Resources
}

func (unit *unit) unitInit(server *Server) {
	unit.server = server
	unit.id = uid.GenUid()
	unit.alive.Store(true)
}

func (unit *unit) GetServer() *Server {
	return unit.server
}

func (unit *unit) GetId() uid.Uid {
	return unit.id
}

func (unit *unit) IsBusy() bool {
	action := unit.GetLastAction()
	return action != nil && !action.Finised.Load()
}

func (unit *unit) GetAlive() bool {
	return unit.alive.Load()
}

func (unit *unit) SetAlive(new bool) {
	unit.alive.Store(new)
}

func (unit *unit) GetPosition() Point {
	return unit.position.Load()
}

func (unit *unit) SetPosition(new_pos Point) {
	unit.position.Store(new_pos)
}

func (unit *unit) GetLastAction() *Action {
	return unit.lastAction.Load()
}

func (unit *unit) ModifyPosition(cb func(Point) Point) (Point, Point) {
	return unit.position.Modify(cb)
}

func (unit *unit) IsUpgraded() bool {
	return unit.upgraded.Load()
}

func (unit *unit) SetUpgraded(new bool) {
	unit.upgraded.Store(new)
}

func (unit *unit) ObserveDistance() int {
	return 0
}

func (unit *unit) GetInventory() model.Resources {
	unit.inventoryLock.RLock()
	defer unit.inventoryLock.RUnlock()
	return unit.inventory
}

func (unit *unit) ModifyInventory(cb func(model.Resources) model.Resources) {
	unit.inventoryLock.Lock()
	defer unit.inventoryLock.Unlock()
	unit.inventory = cb(unit.inventory)
}

func (unit *unit) SetInventory(newInv model.Resources) {
	unit.inventoryLock.Lock()
	defer unit.inventoryLock.Unlock()
	unit.inventory = newInv
}

func startAction(this extendedUnit, action *Action, supported []ActionOpCode) error {
	if this == nil || action == nil {
		panic("cannot work nil")
	}
	if action.Finised.Load() {
		panic("cannot start finished action")
	}

	issupported := false
	for _, op := range supported {
		if op == action.OpCode {
			issupported = true
			break
		}
	}
	if !issupported {
		return UnsuportedActionError{
			Tried:     action.OpCode,
			Supported: supported,
		}
	}

	lastAction := this.GetLastAction()
	if lastAction != nil && !lastAction.Finised.Load() {
		return UnitBusyError{}
	}

	cost := action.OpCode.GetCost(this)
	if this.GetOwner() != uid.ServerUid {
		owner := this.GetServer().GetPlayerFromId(this.GetOwner())
		if owner == nil {
			panic("could not find owner")
		}

		var hadEnough bool
		var had model.Resources
		owner.ModifyResources(func(res model.Resources) model.Resources {
			if res.EnoughFor(cost.Resources) < 1 {
				hadEnough = false
				had = res
				return res
			}
			hadEnough = true
			return res.Sub(cost.Resources)
		})
		if !hadEnough {
			return NotEnoughResourcesError{
				Required:  cost.Resources,
				Available: had,
			}
		}
	}

	this.getUnit().lastAction.Store(action)

	return nil
}

func tick(this IUnit) {
	action := this.GetLastAction()
	if action == nil {
		return
	}

	if action.Finised.Load() {
		return
	}

	costs := action.OpCode.GetCost(this)
	if this.GetServer().Ticker().GetTickNumber()-action.StartedAtTick < costs.Cast {
		return
	}

	// action is finished

	action.Finised.Store(true)

	report := ApplyAction(action, this)
	this.GetServer().AddReport(report)
}
