package entities

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/events"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_server/server"
	. "github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
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
	if asOwner, ok := unit.this.(IOwnerEntity); ok {
		for _, child := range asOwner.CopyEntityList() {
			child.Unregister()
		}
	}

	if !unit.alive.Swap(false) {
		// was already dead
		return
	}

	unit.server.RemoveEntity(unit.id)

	unit.server.Events().Emit(&server.UnitDespawnEvent{
		Unit: unit.this,
		AABB: unit.GetAABB(),
	})
}

func (unit *unit) GetPosition() Point {
	return unit.position.Load()
}

func (unit *unit) SetPosition(new_pos Point) {
	prevValue := unit.position.Store(new_pos)
	unit.movedEvents.Emit(spatialmap.ObjectMovedEvent{
		From: prevValue,
		To:   new_pos,
	})
	if unit.server != nil {
		unit.server.Events().Emit(&server.UnitMovedEvent{
			Unit: unit.this,
			From: prevValue,
			To:   new_pos,
		})
	}
}

func (unit *unit) GetAABB() AABB {
	return AABB{
		From: unit.GetPosition(),
		Size: Point{X: 1, Y: 1},
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
			To:   new,
		})
		if unit.server != nil {
			unit.server.Events().Emit(&server.UnitMovedEvent{
				Unit: unit.this,
				From: old,
				To:   new,
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
		return
	}

	if unit.server == nil {
		return
	}
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

func (unit *unit) startAction(
	action *Action,
	supported []ActionOpCode,
	onFinished func(),
) error {
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
	if player, ok := unit.GetServer().GetEntityOwner(unit.id).(*Player); ok {
		var hadEnough bool
		var had model.Resources
		player.ModifyResources(func(res model.Resources) model.Resources {
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

	action.StartedAtTick = unit.server.Ticker().GetTickNumber()
	unit.lastAction.Store(action)

	unit.server.Events().Emit(&UnitStartedActionEvent{
		Unit:   unit.this,
		Pos:    unit.GetPosition(),
		Action: action,
	})

	go (func() {
		costs := action.OpCode.GetCost(unit.this)

		<-time.After(unit.server.Ticker().TickDuration() * time.Duration(costs.Cast))
		if !unit.GetAlive() {
			return
		}

		action.Finised.Store(true)

		report := ApplyAction(action, unit.this)
		if len(report.GetReport().ReportId) > 0 {
			unit.GetServer().AddReport(report)
		}

		unit.server.Events().Emit(&UnitFinishedActionEvent{
			Unit:   unit.this,
			Pos:    unit.GetPosition(),
			Action: action,
			Report: report,
		})

		if onFinished != nil {
			onFinished()
		}
	})()

	return nil
}

func (unit *unit) Register() {
	unit.server.RegisterEntity(unit.this)
	unit.server.Events().Emit(&server.UnitSpawnEvent{
		Unit: unit.this,
		AABB: unit.GetAABB(),
	})
}

func (unit *unit) Unregister() {
	unit.SetDead()
}
