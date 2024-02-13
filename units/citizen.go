package units

import (
	"sync"

	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
)

type CitizenUnit struct {
	unit
	lock        sync.RWMutex    
	owner       uid.Uid
	lastEatenAt int
}

func NewCitizenUnit(server *Server, owner uid.Uid) *CitizenUnit {
	citizen := new(CitizenUnit)
	citizen.unitInit(server)
	citizen.lastEatenAt = server.Ticker().GetTickNumber()
	citizen.owner = owner

	return citizen
}

func (citizen *CitizenUnit) GetOwner() uid.Uid {
	return citizen.owner
}

func (citizen *CitizenUnit) Tick() {
	citizen.lock.Lock()
	defer citizen.lock.Unlock()

	server := citizen.server
	ticker := server.Ticker()

	feedInterval := server.GetSetup().CitizenFeedingRate

	if ticker.GetTickNumber()-citizen.lastEatenAt > feedInterval {
	}
}
