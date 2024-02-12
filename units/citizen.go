package units

import (
    . "creeps.heav.fr/server"
)

type CitizenUnit struct {
    unit
    owner Uid
    lastEatenAt int
}

func NewCitizenUnit(server *Server, owner Uid) *CitizenUnit {
    citizen := new(CitizenUnit)
    citizen.unitInit(server)
    citizen.lastEatenAt = server.Ticker().GetTickNumber()
    citizen.owner = owner

    server.RegisterUnit(citizen)
    
    return citizen
}

func (citizen *CitizenUnit) GetOwner() Uid {
    return citizen.owner
}

func (citizen *CitizenUnit) Tick() {
    server := citizen.server
    ticker := server.Ticker()

    feedInterval := server.GetSetup().CitizenFeedingRate

    if ticker.GetTickNumber() - citizen.lastEatenAt > feedInterval {
    }
}
