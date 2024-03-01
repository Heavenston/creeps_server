package model

import (
	"github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

type StatusResponse struct {
	Running bool `json:"running"`
}

type Player struct {
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Units        int       `json:"units"`
	Buildings    int       `json:"buildings"`
	Resources    Resources `json:"resources"`
	Achievements []string  `json:"achievements"`
}

type StatisticsResponse struct {
	ServerId    string     `json:"serverId"`
	GameRunning bool       `json:"gameRunning"`
	Tick        int        `json:"tick"`
	Dimension   geom.Point `json:"dimension"`
	Players     []Player   `json:"players"`
}

type CostResponse struct {
	Resources
	Cast int `json:"cast"`
}

type CostsResponse struct {
	BuildHousehold   CostResponse `json:"buildHousehold"`
	BuildRoad        CostResponse `json:"buildRoad"`
	BuildSawmill     CostResponse `json:"buildSawmill"`
	BuildSmeltery    CostResponse `json:"buildSmeltery"`
	BuildTownHall    CostResponse `json:"buildTownHall"`
	Dismantle        CostResponse `json:"dismantle"`
	Farm             CostResponse `json:"farm"`
	FetchMessage     CostResponse `json:"fetchMessage"`
	FireBomberBot    CostResponse `json:"fireBomberBot"`
	FireTurret       CostResponse `json:"fireTurret"`
	Gather           CostResponse `json:"gather"`
	Move             CostResponse `json:"move"`
	Noop             CostResponse `json:"noop"`
	Observe          CostResponse `json:"observe"`
	RefineCopper     CostResponse `json:"refineCopper"`
	RefineWoodPlank  CostResponse `json:"refineWoodPlank"`
	SendMessage      CostResponse `json:"sendMessage"`
	SpawnBomberBot   CostResponse `json:"spawnBomberBot"`
	SpawnTurret      CostResponse `json:"spawnTurret"`
	Unload           CostResponse `json:"unload"`
	UpgradeBomberBot CostResponse `json:"upgradeBomberBot"`
	UpgradeCitizen   CostResponse `json:"upgradeCitizen"`
	UpgradeTurret    CostResponse `json:"upgradeTurret"`
}

type SetupResponse struct {
	CitizenFeedingRate int  `json:"citizenFeedingRate"`
	EnableGC           bool `json:"enableGC"`
	GcTickRate         int  `json:"gcTickRate"`
	EnableEnemies      bool `json:"enableEnemies"`
	EnemyTickRate      int  `json:"enemyTickRate"`
	// disabled from json for epita compitibility
	EnemyBaseTickRate  int        `json:"-"`
	FoodGatherRate     int        `json:"foodGatherRate"`
	MaxLoad            int        `json:"maxLoad"`
	MaxMissesPerPlayer int        `json:"maxMissesPerPlayer"`
	MaxMissesPerUnit   int        `json:"maxMissesPerUnit"`
	OilGatherRate      int        `json:"oilGatherRate"`
	RockGatherRate     int        `json:"rockGatherRate"`
	ServerId           string     `json:"serverId"`
	TicksPerSecond     float64    `json:"ticksPerSecond"`
	TrackAchievements  bool       `json:"trackAchievements"`
	WoodGatherRate     int        `json:"woodGatherRate"`
	WorldDimension     geom.Point `json:"worldDimension"`
}

type InitResponse struct {
	Citizen1Id           *uid.Uid       `json:"citizen1Id"`
	Citizen2Id           *uid.Uid       `json:"citizen2Id"`
	Costs                *CostsResponse `json:"costs"`
	Error                *string        `json:"error"`
	HouseholdCoordinates *geom.Point    `json:"householdCoordinates"`
	Login                string         `json:"login"`
	PlayerId             *int16         `json:"playerId"`
	Resources            *Resources     `json:"resources"`
	Setup                *SetupResponse `json:"setup"`
	Tick                 int            `json:"tick"`
	TownHallCoordinates  *geom.Point    `json:"townHallCoordinates"`
}

type CommandResponse struct {
	OpCode    ActionOpCode `json:"opcode"`
	ReportId  *uid.Uid     `json:"reportId"`
	ErrorCode *string      `json:"errorCode"`
	Error     *string      `json:"error"`
	Login     string       `json:"login"`
	UnitId    *uid.Uid     `json:"unitId"`
	Misses    int          `json:"misses"`
}
