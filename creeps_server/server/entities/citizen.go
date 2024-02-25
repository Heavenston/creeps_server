package entities

import (
	"sync"

	"github.com/heavenston/creeps_server/creeps_lib/model"
	. "github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
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
	citizen.this = citizen
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

func (citizen *CitizenUnit) StartAction(action *Action, onFinished func()) error {
	err := citizen.startAction(action, []ActionOpCode{
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
	}, onFinished)
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
	player, ok := server.GetEntity(citizen.GetOwner()).(*Player)
	if !ok {
		log.Warn().
			Str("citizen_id", string(citizen.id)).
			Str("player_id", string(citizen.owner)).
			Msg("[CITIZEN] Could not find owner player (code kms initiated)")
		citizen.SetDead()
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
			citizen.SetDead()
		}
	}
}
