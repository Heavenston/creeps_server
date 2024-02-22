package units

import (
	"sync"
	"sync/atomic"

	"creeps.heav.fr/epita_api/model"
	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/server"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/spatialmap"
	"creeps.heav.fr/uid"
)

type extendedUnit interface {
	IUnit
	getUnit() *unit
}

// See server.go's IUnit interface to explain its functions
type unit struct {
	this IUnit
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
	movedEvents   events.EventProvider[spatialmap.ObjectMovedEvent]
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

func (unit *unit) SetDead() {
	if !unit.alive.Swap(false) {
		// was already dead
		return
	}

	if unit.server == nil {
		return
	}

	unit.server.RemoveUnit(unit.id)
}

func (unit *unit) GetPosition() Point {
	return unit.position.Load()
}

func (unit *unit) SetPosition(new_pos Point) {
	prevValue := unit.position.Store(new_pos)
	unit.movedEvents.Emit(spatialmap.ObjectMovedEvent{
		From: prevValue,
		To: new_pos,
	})
	if unit.server != nil {
		unit.server.Events().Emit(&server.UnitMovedEvent {
			Unit: unit.this,
			From: prevValue,
			To: new_pos,
		})
	}
}

func (unit *unit) GetAABB() AABB {
	return AABB {
		From: unit.GetPosition(),
		Size: Point { X: 1, Y: 1 },
	}
}

func (unit *unit) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return &unit.movedEvents
}

func (unit *unit) GetLastAction() *Action {
	return unit.lastAction.Load()
}

func (unit *unit) ModifyPosition(cb func(Point) Point) (Point, Point) {
	old, new := unit.position.Modify(cb)
	if old != new {
		unit.movedEvents.Emit(spatialmap.ObjectMovedEvent{
			From: old,
			To: new,
		})
		if unit.server != nil {
			unit.server.Events().Emit(&server.UnitMovedEvent {
				Unit: unit.this,
				From: old,
				To: new,
			})
		}
	}
	return old, new
}

func (unit *unit) IsUpgraded() bool {
	return unit.upgraded.Load()
}

func (unit *unit) SetUpgraded() {
	was := unit.upgraded.Swap(true)
	if was {
		return;
	}

	if unit.server == nil {
		return
	}

	unit.server.Events().Emit(&server.UnitUpgradedEvent{
		Unit: unit.this,
	})
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

func (unit *unit) startAction(action *Action, supported []ActionOpCode) error {
	if unit == nil || action == nil {
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

	lastAction := unit.GetLastAction()
	if lastAction != nil && !lastAction.Finised.Load() {
		return UnitBusyError{}
	}

	cost := action.OpCode.GetCost(unit.this)
	if unit.this.GetOwner() != uid.ServerUid {
		owner := unit.GetServer().GetPlayerFromId(unit.this.GetOwner())
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

	unit.lastAction.Store(action)

	return nil
}

func (unit *unit) tick() {
	action := unit.GetLastAction()
	if action == nil {
		return
	}

	if action.Finised.Load() {
		return
	}

	costs := action.OpCode.GetCost(unit.this)
	if unit.GetServer().Ticker().GetTickNumber()-action.StartedAtTick < costs.Cast {
		return
	}

	// action is finished

	action.Finised.Store(true)

	report := ApplyAction(action, unit.this)
	unit.GetServer().AddReport(report)
}
