package units

import (
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
)

// See server.go's IUnit interface to explain its functions
type unit struct {
	server     *Server
	id         uid.Uid
	alive      bool
	position   Point
	lastAction *Action
}

func (unit *unit) unitInit(server *Server) {
	unit.server = server
	unit.id = uid.GenUid()
	unit.alive = true
}

func (unit *unit) GetServer() *Server {
	return unit.server
}

func (unit *unit) GetId() uid.Uid {
	return unit.id
}

func (unit *unit) GetAlive() bool {
	return unit.alive
}

func (unit *unit) SetAlive(new bool) {
	unit.alive = new
}

func (unit *unit) GetPosition() Point {
	return unit.position
}

func (unit *unit) SetPosition(new_pos Point) {
	unit.position = new_pos
}

func (unit *unit) GetLastAction() *Action {
	return unit.lastAction
}

func (unit *unit) SetLastAction(action *Action) {
	unit.lastAction = action
}
