package units

import (
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
)

type CitizenUnit struct {
	unit
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
	server := citizen.server
	ticker := server.Ticker()

	feedInterval := server.GetSetup().CitizenFeedingRate

	if ticker.GetTickNumber()-citizen.lastEatenAt > feedInterval {
	}
}
