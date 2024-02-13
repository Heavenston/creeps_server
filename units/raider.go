package units

import (
	"sync"

	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
)

type RaiderUnit struct {
	unit
	lock   sync.RWMutex    
	target AtomicPoint
}

func NewRaiderUnit(server *Server, target Point) *RaiderUnit {
	raider := new(RaiderUnit)
	raider.unitInit(server)
	raider.target.Store(target)
	return raider
}

func (raider *RaiderUnit) GetOwner() uid.Uid {
	return uid.ServerUid
}

func (raider *RaiderUnit) GetTarget() Point {
	return raider.target.Load()
}

// returns the previous value
func (raider *RaiderUnit) SetTarget(target Point) Point {
	return raider.target.Store(target)
}

func (raider *RaiderUnit) Tick() {
	raider.lock.Lock()
	defer raider.lock.Unlock()
}
