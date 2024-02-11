package server

import "creeps.heav.fr/geom"

type iunit interface {
    setId(id int)
    GetId() int
    GetOwner() *UnitOwner
    GetPosition() *geom.Point
    GetLastAction() *Action
    SetLastAction(action *Action)
}

type Unit struct {
    id int
    position geom.Point
    lastAction *Action
}

func (unit *Unit) GetId() int {
    return unit.id
}

func (unit *Unit) setId(id int) {
    unit.id = id
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
    target *geom.Point
}

func NewRaiderUnit() *RaiderUnit {
    return new(RaiderUnit)
}

func (raider *RaiderUnit) GetOwner() *UnitOwner {
    return nil
}
