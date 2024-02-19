package units

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
	"github.com/rs/zerolog/log"
)

type CitizenUnit struct {
	unit
	lock        sync.Mutex
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

func (citizen *CitizenUnit) GetOpCode() string {
	return "citizen"
}

func (citizen *CitizenUnit) GetUpgradeCosts() *model.CostResponse {
	return &citizen.GetServer().GetCosts().UpgradeCitizen
}

func (citizen *CitizenUnit) GetOwner() uid.Uid {
	return citizen.owner
}

func (citizen *CitizenUnit) ObserveDistance() int {
	if citizen.IsUpgraded() {
		// FIXME: Is it correct ?
		return 7
	}
	return 6
}

func (citizen *CitizenUnit) StartAction(action *Action) error {
	err := startAction(citizen, action, []ActionOpCode{
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
		OpCodeFarm,
		OpCodeUnload,
		OpCodeRefineCopper,
		OpCodeRefineWoodPlank,

		OpCodeObserve,
	})
	if err != nil {
		return err
	}
	return nil
}

func (citizen *CitizenUnit) Tick() {
	citizen.lock.Lock()
	defer citizen.lock.Unlock()

	tick(citizen)

	server := citizen.server
	ticker := server.Ticker()
	player := server.GetPlayerFromId(citizen.GetOwner())
	if player == nil {
		log.Error().
			Msg("[CITIZEN] Could not find owner player (code kms initiated)")
		citizen.SetAlive(false)
		return
	}

	feedInterval := server.GetSetup().CitizenFeedingRate
	feedAmount := 1
	if citizen.IsUpgraded() {
		feedAmount = 2
	}

	if ticker.GetTickNumber()-citizen.lastEatenAt > feedInterval {
		var couldFeed bool
		player.ModifyResources(func(res model.Resources) model.Resources {
			if res.Food < feedAmount {
				couldFeed = false
				return res
			}
			couldFeed = true
			res.Food -= feedAmount
			return res
		})
		if couldFeed {
			citizen.lastEatenAt = ticker.GetTickNumber()
		} else {
			citizen.SetAlive(false)
		}
	}
}
