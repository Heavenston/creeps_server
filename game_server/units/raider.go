package units

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
	"github.com/rs/zerolog/log"
)

type RaiderUnit struct {
	unit
	lock   sync.RWMutex
	target Point
}

func NewRaiderUnit(server *Server, target Point) *RaiderUnit {
	raider := new(RaiderUnit)
	raider.unitInit(server)
	raider.target = target
	raider.this = raider
	return raider
}

// for the extendedUnit interface
func (raider *RaiderUnit) getUnit() *unit {
	return &raider.unit
}

func (raider *RaiderUnit) GetOpCode() string {
	return "raider"
}

func (raider *RaiderUnit) GetUpgradeCosts() *model.CostResponse {
	return nil
}

func (raider *RaiderUnit) GetOwner() uid.Uid {
	return uid.ServerUid
}

func (raider *RaiderUnit) GetTarget() Point {
	return raider.target
}

func (raider *RaiderUnit) StartAction(action *Action) error {
	err := raider.startAction(action, []ActionOpCode {
		OpCodeMoveDown,
		OpCodeMoveUp,
		OpCodeMoveLeft,
		OpCodeMoveRight,
	})
	if err != nil {
		return err
	}
	return nil
}

func (raider *RaiderUnit) Tick() {
	raider.lock.Lock()
	defer raider.lock.Unlock()

	raider.tick()

	position := raider.GetPosition()

	// busy = do nothing
	if action := raider.GetLastAction(); action != nil && !action.Finised.Load() {
		return
	}

	if raider.target == position {
		raider.SetAlive(false)
		return
	}

	diff := raider.target.Sub(position)
	newAction := new(Action)
	newAction.StartedAtTick = raider.GetServer().Ticker().GetTickNumber()

	if mathutils.AbsInt(diff.X) > mathutils.AbsInt(diff.Y) {
		if diff.X < 0 {
			newAction.OpCode = OpCodeMoveRight
		} else {
			newAction.OpCode = OpCodeMoveLeft
		}
	} else {
		if diff.Y < 0 {
			newAction.OpCode = OpCodeMoveDown
		} else {
			newAction.OpCode = OpCodeMoveUp
		}
	}

	err := raider.StartAction(newAction)
	if err != nil {
		log.Error().
			Any("action", newAction).
			Err(err).
			Msg("[RAIDER] Could not start action")
	}
}
