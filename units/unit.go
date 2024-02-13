package units

import (
	"sync/atomic"

	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
)

// See server.go's IUnit interface to explain its functions
type unit struct {
	// read-only (no lock)
	server *Server
	// read-only (no lock)
	id         uid.Uid
	alive      atomic.Bool
	position   AtomicPoint
	lastAction atomic.Pointer[Action]
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

func (unit *unit) ModifyPosition(cb func (Point) Point) (Point, Point) {
	return unit.position.Modify(cb)
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

	action.OpCode.ApplyOn(this)
}
