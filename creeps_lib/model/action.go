package model

import (
	"reflect"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
)

type ActionOpCode string

const (
	OpCodeMoveLeft  ActionOpCode = "move:left"
	OpCodeMoveRight              = "move:right"
	OpCodeMoveUp                 = "move:up"
	OpCodeMoveDown               = "move:down"

	OpCodeObserve   = "observe"
	OpCodeGather    = "gather"
	OpCodeUnload    = "unload"
	OpCodeFarm      = "farm"
	OpCodeDismantle = "dismantle"
	OpCodeUpgrade   = "upgrade"

	OpCodeRefineCopper    = "refine:copper"
	OpCodeRefineWoodPlank = "refine:wood-plank"

	OpCodeBuildTownHall  = "build:town-hall"
	OpCodeBuildHousehold = "build:household"
	OpCodeBuildSawmill   = "build:sawmill"
	OpCodeBuildSmeltery  = "build:smeltery "
	OpCodeBuildRoad      = "build:road"

	OpCodeSpawnTurret    = "spawn:turret"
	OpCodeSpawnBomberBot = "spawn:bomber-bot"

	OpCodeFireTurret    = "fire:turret"
	OpCodeFireBomberBot = "fire:bomber-bot"
)

func (opcode ActionOpCode) IsValid() bool {
	switch opcode {
	case OpCodeMoveLeft:
		return true
	case OpCodeMoveRight:
		return true
	case OpCodeMoveUp:
		return true
	case OpCodeMoveDown:
		return true
	case OpCodeObserve:
		return true
	case OpCodeGather:
		return true
	case OpCodeUnload:
		return true
	case OpCodeFarm:
		return true
	case OpCodeDismantle:
		return true
	case OpCodeUpgrade:
		return true
	case OpCodeRefineCopper:
		return true
	case OpCodeRefineWoodPlank:
		return true
	case OpCodeBuildTownHall:
		return true
	case OpCodeBuildHousehold:
		return true
	case OpCodeBuildSawmill:
		return true
	case OpCodeBuildSmeltery:
		return true
	case OpCodeBuildRoad:
		return true
	case OpCodeSpawnTurret:
		return true
	case OpCodeSpawnBomberBot:
		return true
	case OpCodeFireTurret:
		return true
	case OpCodeFireBomberBot:
		return true
	}
	return false
}

func (opcode ActionOpCode) GetReportType() reflect.Type {
	switch opcode {
	case OpCodeMoveLeft:
		fallthrough
	case OpCodeMoveRight:
		fallthrough
	case OpCodeMoveUp:
		fallthrough
	case OpCodeMoveDown:
		return reflect.TypeFor[MoveReport]()

	case OpCodeObserve:
		return reflect.TypeFor[ObserveReport]()
	case OpCodeGather:
		return reflect.TypeFor[GatherReport]()
	case OpCodeUnload:
		return reflect.TypeFor[UnloadReport]()
	case OpCodeFarm:
		return reflect.TypeFor[FarmReport]()
	case OpCodeDismantle:
		return reflect.TypeFor[DismantleReport]()
	case OpCodeUpgrade:
		return reflect.TypeFor[UpgradeReport]()

	case OpCodeRefineCopper:
		fallthrough
	case OpCodeRefineWoodPlank:
		return reflect.TypeFor[RefineReport]()

	case OpCodeBuildTownHall:
		fallthrough
	case OpCodeBuildHousehold:
		fallthrough
	case OpCodeBuildSawmill:
		fallthrough
	case OpCodeBuildSmeltery:
		fallthrough
	case OpCodeBuildRoad:
		return reflect.TypeFor[BuildReport]()

	case OpCodeSpawnTurret:
		fallthrough
	case OpCodeSpawnBomberBot:
		return reflect.TypeFor[SpawnReport]()

	case OpCodeFireTurret:
		fallthrough
	case OpCodeFireBomberBot:
		return reflect.TypeFor[FireReport]()
	}
	panic("invalid opcode")
}

func (opcode ActionOpCode) GetCost(costs *CostsResponse, upgradeCost *CostResponse) *CostResponse {
	switch opcode {
	case OpCodeMoveLeft:
		return &costs.Move
	case OpCodeMoveRight:
		return &costs.Move
	case OpCodeMoveUp:
		return &costs.Move
	case OpCodeMoveDown:
		return &costs.Move
	case OpCodeObserve:
		return &costs.Observe
	case OpCodeGather:
		return &costs.Gather
	case OpCodeUnload:
		return &costs.Unload
	case OpCodeFarm:
		return &costs.Farm
	case OpCodeDismantle:
		return &costs.Dismantle
	case OpCodeUpgrade:
		return upgradeCost
	case OpCodeRefineCopper:
		return &costs.RefineCopper
	case OpCodeRefineWoodPlank:
		return &costs.RefineWoodPlank
	case OpCodeBuildTownHall:
		return &costs.BuildTownHall
	case OpCodeBuildHousehold:
		return &costs.BuildHousehold
	case OpCodeBuildSawmill:
		return &costs.BuildSawmill
	case OpCodeBuildSmeltery:
		return &costs.BuildSmeltery
	case OpCodeBuildRoad:
		return &costs.BuildRoad
	case OpCodeSpawnTurret:
		return &costs.SpawnTurret
	case OpCodeSpawnBomberBot:
		return &costs.SpawnBomberBot
	case OpCodeFireTurret:
		return &costs.FireTurret
	case OpCodeFireBomberBot:
		return &costs.FireBomberBot
	}
	panic("invalid opcode")
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

func (opcode ActionOpCode) ParameterType() reflect.Type {
	switch opcode {
	case OpCodeFireTurret:
		return reflect.TypeFor[FireParameter]()
	case OpCodeFireBomberBot:
		return reflect.TypeFor[FireParameter]()
	default:
		return nil
	}
}
