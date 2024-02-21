package units

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
)

type BomberBotUnit struct {
	unit
	lock        sync.Mutex
	owner       uid.Uid
}

func NewBomberBotUnit(server *Server, owner uid.Uid) *BomberBotUnit {
	bomberBot := new(BomberBotUnit)
	bomberBot.unitInit(server)
	bomberBot.owner = owner
	bomberBot.this = bomberBot
	return bomberBot
}

// for the extendedUnit interface
func (bomberBot *BomberBotUnit) getUnit() *unit {
	return &bomberBot.unit
}

func (bomberBot *BomberBotUnit) GetOpCode() string {
	return "bomberBot"
}

func (bomberBot *BomberBotUnit) GetUpgradeCosts() *model.CostResponse {
	return &bomberBot.GetServer().GetCosts().UpgradeBomberBot
}

func (bomberBot *BomberBotUnit) GetOwner() uid.Uid {
	return bomberBot.owner
}

func (bomberBot *BomberBotUnit) ObserveDistance() int {
	if bomberBot.IsUpgraded() {
		// FIXME: Is it correct ?
		return 7
	}
	return 5
}

func (bomberBot *BomberBotUnit) StartAction(action *Action) error {
	err := bomberBot.startAction(action, []ActionOpCode{
		OpCodeUpgrade,
		OpCodeFireBomberBot,
	})
	if err != nil {
		return err
	}
	return nil
}

func (bomberBot *BomberBotUnit) Tick() {
	bomberBot.lock.Lock()
	defer bomberBot.lock.Unlock()

	bomberBot.tick()
}
