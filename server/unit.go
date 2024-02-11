package server

import "creeps.heav.fr/geom"

type IUnit interface {
    GetId() Uid
    // the id of the owner, note: can be the server by way of ServerUid
    GetOwner() Uid
    GetPosition() *geom.Point
    GetLastAction() *Action
    SetLastAction(action *Action)
    Tick()
}

type Unit struct {
    server *Server
    id Uid
    position geom.Point
    lastAction *Action
}

func (unit *Unit) unitInit(server *Server) {
    unit.id = GenUid()
}

func (unit *Unit) GetId() Uid {
    return unit.id
}

func (unit *Unit) GetPosition() *geom.Point {
    return &unit.position
}

func (unit *Unit) GetLastAction() *Action {
    return unit.lastAction
}

func (unit *Unit) SetLastAction(action *Action) {
    unit.lastAction = action
}

type RaiderUnit struct {
    Unit
    target geom.Point
}

func NewRaiderUnit(server *Server, target geom.Point) *RaiderUnit {
    raider := new(RaiderUnit)
    raider.unitInit(server)
    raider.target = target
    return raider
}

func (raider *RaiderUnit) GetOwner() Uid {
    return ServerUid
}

func (raider *RaiderUnit) Tick() {
}

type CitizenUnit struct {
    Unit
    owner Uid
    lastEatenAt int
}

func NewCitizenUnit(server *Server, owner Uid) *CitizenUnit {
    citizen := new(CitizenUnit)
    citizen.unitInit(server)
    citizen.lastEatenAt = server.ticker.tickNumber
    citizen.owner = owner

    server.RegisterUnit(citizen)
    
    return citizen
}

func (citizen *CitizenUnit) GetOwner() Uid {
    return citizen.owner
}

func (raider *CitizenUnit) Tick() {
}
