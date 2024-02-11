package units

import (
    . "creeps.heav.fr/server"
    . "creeps.heav.fr/geom"
)

type RaiderUnit struct {
    unit
    target Point
}

func NewRaiderUnit(server *Server, target Point) *RaiderUnit {
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

