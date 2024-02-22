package entities

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
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
	return 6
}

func (turret *TurretUnit) StartAction(action *Action) error {
	err := turret.startAction(action, []ActionOpCode{
		OpCodeUpgrade,
		OpCodeFireTurret,
	})
	if err != nil {
		return err
	}
	return nil
}

func (turret *TurretUnit) Tick() {
	turret.lock.Lock()
	defer turret.lock.Unlock()

	turret.tick()
}

