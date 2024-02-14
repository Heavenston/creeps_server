package units

import (
	"sync"

	"creeps.heav.fr/api/model"
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

// for the extendedUnit interface
func (citizen *CitizenUnit) getUnit() *unit {
	return &citizen.unit
}

func (citizen *CitizenUnit) GetUpgradeCosts() *model.CostResponse {
	return &citizen.GetServer().GetCosts().UpgradeCitizen
}

func (citizen *CitizenUnit) GetOwner() uid.Uid {
	return citizen.owner
}

func (citizen *CitizenUnit) StartAction(action *Action) error {
	err := startAction(citizen, action, []ActionOpCode {
		OpCodeMoveDown,
		OpCodeMoveUp,
		OpCodeMoveLeft,
		OpCodeMoveRight,

		OpCodeBuildHousehold,
		OpCodeBuildRoad,
		OpCodeBuildSawmill,
		OpCodeBuildSmeltery,
		OpCodeBuildTownHall,

		OpCodeUpgrade,
		OpCodeSpawnBomberBot,
		OpCodeSpawnTurret,

		OpCodeGather,
		OpCodeRefineCopper,
		OpCodeRefineWoodPlank,
	})
	if err != nil {
		return err
	}
	return nil
}

func (citizen *CitizenUnit) Tick() {
	citizen.lock.Lock()
	defer citizen.lock.Unlock()

	server := citizen.server
	ticker := server.Ticker()

	feedInterval := server.GetSetup().CitizenFeedingRate

	if ticker.GetTickNumber()-citizen.lastEatenAt > feedInterval {
		citizen.lastEatenAt = ticker.GetTickNumber()
	}
}
