package server

import (
	"sync/atomic"

	"creeps.heav.fr/api/model"
	"creeps.heav.fr/uid"
	. "creeps.heav.fr/geom"
)

type ActionOpCode string

const (
	OpCodeMoveLeft        ActionOpCode = "move:left"
	OpCodeMoveRight       ActionOpCode = "move:right"
	OpCodeMoveUp          ActionOpCode = "move:up"
	OpCodeMoveDown        ActionOpCode = "move:down"
	OpCodeObserve         ActionOpCode = "observe"
	OpCodeGather          ActionOpCode = "gather"
	OpCodeDismantle       ActionOpCode = "dismantle"
	OpCodeUpgrade         ActionOpCode = "upgrade"
	OpCodeRefineCopper    ActionOpCode = "refine:copper"
	OpCodeRefineWoodPlank ActionOpCode = "refine:wood-plank"
	OpCodeBuildTownHall   ActionOpCode = "build:town-hall"
	OpCodeBuildHousehold  ActionOpCode = "build:household"
	OpCodeBuildSawmill    ActionOpCode = "build:sawmill"
	OpCodeBuildSmeltery   ActionOpCode = "build:smeltery "
	OpCodeBuildRoad       ActionOpCode = "build:road"
	OpCodeSpawnTurret     ActionOpCode = "spawn:turret"
	OpCodeSpawnBomberBot  ActionOpCode = "spawn:bomber-bot"
	OpCodeFireTurret      ActionOpCode = "fire:turret"
	OpCodeFireBomberBot   ActionOpCode = "fire:bomber-bot"
)

func (opcode ActionOpCode) GetCost(unit IUnit) *model.CostResponse {
	srv := unit.GetServer()
	switch opcode {
	case OpCodeMoveLeft:
		return &srv.costs.Move
	case OpCodeMoveRight:
		return &srv.costs.Move
	case OpCodeMoveUp:
		return &srv.costs.Move
	case OpCodeMoveDown:
		return &srv.costs.Move
	case OpCodeObserve:
		return &srv.costs.Observe
	case OpCodeGather:
		return &srv.costs.Gather
	case OpCodeDismantle:
		return &srv.costs.Dismantle
	case OpCodeUpgrade:
		return unit.GetUpgradeCosts()
	case OpCodeRefineCopper:
		return &srv.costs.RefineCopper
	case OpCodeRefineWoodPlank:
		return &srv.costs.RefineWoodPlank
	case OpCodeBuildTownHall:
		return &srv.costs.BuildTownHall
	case OpCodeBuildHousehold:
		return &srv.costs.BuildHousehold
	case OpCodeBuildSawmill:
		return &srv.costs.BuildSawmill
	case OpCodeBuildSmeltery:
		return &srv.costs.BuildSmeltery
	case OpCodeBuildRoad:
		return &srv.costs.BuildRoad
	case OpCodeSpawnTurret:
		return &srv.costs.SpawnTurret
	case OpCodeSpawnBomberBot:
		return &srv.costs.SpawnBomberBot
	case OpCodeFireTurret:
		return &srv.costs.FireTurret
	case OpCodeFireBomberBot:
		return &srv.costs.FireBomberBot
	}
	return nil
}

func (opcode ActionOpCode) MoveDirection() Point {
	switch opcode {
	case OpCodeMoveDown:
		return Point{X:0, Y:-1}
	case OpCodeMoveUp:
		return Point{X:0, Y:1}
	case OpCodeMoveLeft:
		return Point{X:-1, Y:0}
	case OpCodeMoveRight:
		return Point{X:1, Y:0}
	default:
		return Point{}
	}
}

// called by unit in units/unit.go when the action is finished
func (opcode ActionOpCode) ApplyOn(unit IUnit) {
	// srv := unit.GetServer()

	movement := opcode.MoveDirection()
	if movement != (Point{}) {
		unit.ModifyPosition(func (pos Point) Point {
			return pos.Add(movement)
		})
	}
	
	switch opcode {
	case OpCodeMoveLeft:
	case OpCodeMoveRight:
	case OpCodeMoveUp:
	case OpCodeMoveDown:

	case OpCodeObserve:
	case OpCodeGather:
	case OpCodeDismantle:
	case OpCodeUpgrade:
	case OpCodeRefineCopper:
	case OpCodeRefineWoodPlank:
	case OpCodeBuildTownHall:
	case OpCodeBuildHousehold:
	case OpCodeBuildSawmill:
	case OpCodeBuildSmeltery:
	case OpCodeBuildRoad:
	case OpCodeSpawnTurret:
	case OpCodeSpawnBomberBot:
	case OpCodeFireTurret:
	case OpCodeFireBomberBot:
	}
}

// every value is read only except
type Action struct {
	OpCode        ActionOpCode
	StartedAtTick int
	ReportId      uid.Uid
	Finised       atomic.Bool
}
