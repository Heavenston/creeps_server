package units

import (
	"errors"
	"sync"

	"creeps.heav.fr/api/model"
	. "creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/uid"
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
	return raider
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
	if action.Finised.Load() {
		panic("cannot start finished action")
	}

	if action.OpCode.MoveDirection() == (Point{}) {
		return UnsuportedActionError{Tried: action.OpCode}
	}

	lastAction := raider.GetLastAction()
	if lastAction != nil || !lastAction.Finised.Load() {
		return UnitBusyError{}
	}

	return nil
}

func (raider *RaiderUnit) Tick() {
	raider.lock.Lock()
	defer raider.lock.Unlock()

	tick(raider)

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
	newAction.ReportId = uid.GenUid()
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
	errors.Unwrap(err)
}
