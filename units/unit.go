package units

import (
    . "creeps.heav.fr/server"
    . "creeps.heav.fr/geom"
)

type unit struct {
    server *Server
    id Uid
    position Point
    lastAction *Action
}

func (unit *unit) unitInit(server *Server) {
    unit.id = GenUid()
}

func (unit *unit) GetServer() *Server {
    return unit.server
}

func (unit *unit) GetId() Uid {
    return unit.id
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

