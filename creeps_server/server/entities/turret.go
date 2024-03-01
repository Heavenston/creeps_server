package entities

import (
	"sync"

	"github.com/heavenston/creeps_server/creeps_lib/model"
	. "github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

type TurretUnit struct {
	unit
	lock        sync.Mutex
	owner       uid.Uid
}

func NewTurretUnit(server *Server, owner uid.Uid) *TurretUnit {
	turret := new(TurretUnit)
	turret.unitInit(server)
	turret.owner = owner
	turret.this = turret
	return turret
}

// for the extendedUnit interface
func (turret *TurretUnit) getUnit() *unit {
	return &turret.unit
}

func (turret *TurretUnit) GetOpCode() string {
	return "turret"
}

func (turret *TurretUnit) GetUpgradeCosts() *model.CostResponse {
	return &turret.GetServer().GetCosts().UpgradeTurret
}

func (turret *TurretUnit) GetOwner() uid.Uid {
	return turret.owner
}

func (turret *TurretUnit) ObserveDistance() int {
	if turret.IsUpgraded() {
		// FIXME: Is it correct ?
		return 7
	}
	return 5
}

func (turret *TurretUnit) StartAction(action *Action, onFinished func()) error {
	err := turret.startAction(action, []model.ActionOpCode{
		model.OpCodeUpgrade,
		model.OpCodeFireTurret,
	}, onFinished)
	if err != nil {
		return err
	}
	return nil
}

func (turret *TurretUnit) Tick() { }

