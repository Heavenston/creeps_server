package server

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

type Action struct {
	OpCode        ActionOpCode
	StartedAtTick int
	ReportId      Uid
}
