package server

import (
	"sync/atomic"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/uid"
)

type ActionOpCode string

const (
	OpCodeMoveLeft ActionOpCode = "move:left"

	OpCodeMoveRight       = "move:right"
	OpCodeMoveUp          = "move:up"
	OpCodeMoveDown        = "move:down"
	OpCodeObserve         = "observe"
	OpCodeGather          = "gather"
	OpCodeUnload          = "unload"
	OpCodeFarm            = "farm"
	OpCodeDismantle       = "dismantle"
	OpCodeUpgrade         = "upgrade"
	OpCodeRefineCopper    = "refine:copper"
	OpCodeRefineWoodPlank = "refine:wood-plank"
	OpCodeBuildTownHall   = "build:town-hall"
	OpCodeBuildHousehold  = "build:household"
	OpCodeBuildSawmill    = "build:sawmill"
	OpCodeBuildSmeltery   = "build:smeltery "
	OpCodeBuildRoad       = "build:road"
	OpCodeSpawnTurret     = "spawn:turret"
	OpCodeSpawnBomberBot  = "spawn:bomber-bot"
	OpCodeFireTurret      = "fire:turret"
	OpCodeFireBomberBot   = "fire:bomber-bot"
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
	case OpCodeUnload:
		return &srv.costs.Unload
	case OpCodeFarm:
		return &srv.costs.Farm
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
		return Point{X: 0, Y: -1}
	case OpCodeMoveUp:
		return Point{X: 0, Y: 1}
	case OpCodeMoveLeft:
		return Point{X: -1, Y: 0}
	case OpCodeMoveRight:
		return Point{X: 1, Y: 0}
	default:
		return Point{}
	}
}

// every value is read only except
type Action struct {
	OpCode        ActionOpCode
	StartedAtTick int
	ReportId      uid.Uid
	Finised       atomic.Bool
}
