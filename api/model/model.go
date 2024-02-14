package model

import (
	"creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
	"creeps.heav.fr/uid"
)

type ResourceKind string

const (
	Copper    ResourceKind = "copper"
	Food                   = "food"
	Oil                    = "oil"
	Rock                   = "rock"
	Wood                   = "wood"
	WoodPlank              = "woodPlank"
)

type Resources struct {
	Rock      int `json:"rock"`
	Wood      int `json:"wood"`
	Food      int `json:"food"`
	Oil       int `json:"oil"`
	Copper    int `json:"copper"`
	WoodPlank int `json:"woodPlank"`
}

func (res *Resources) OfKind(kind ResourceKind) *int {
	switch kind {
	case Rock:
		return &res.Rock
	case Wood:
		return &res.Wood
	case Food:
		return &res.Food
	case Oil:
		return &res.Oil
	case Copper:
		return &res.Copper
	case WoodPlank:
		return &res.WoodPlank
	}
	return nil
}

// return how many times this resources have the other one
// (4 copper 2 rock for 1 copper 1 rock returns 2)
func (res Resources) EnoughFor(other Resources) float64 {
	return mathutils.Min(
		float64(res.Rock)/float64(other.Rock),
		float64(res.Wood)/float64(other.Wood),
		float64(res.Food)/float64(other.Food),
		float64(res.Oil)/float64(other.Oil),
		float64(res.Copper)/float64(other.Copper),
		float64(res.WoodPlank)/float64(other.WoodPlank),
	)
}

func (res *Resources) Remove(other Resources) {
	res.Rock -= other.Rock
	res.Wood -= other.Wood
	res.Food -= other.Food
	res.Oil -= other.Oil
	res.Copper -= other.Copper
	res.WoodPlank -= other.WoodPlank
}

func (res Resources) Sub(other Resources) Resources {
	res.Remove(other)
	return res
}

func (res *Resources) Add(other Resources) {
	res.Rock += other.Rock
	res.Wood += other.Wood
	res.Food += other.Food
	res.Oil += other.Oil
	res.Copper += other.Copper
	res.WoodPlank += other.WoodPlank
}

func (res Resources) Sum(other Resources) Resources {
	res.Add(other)
	return res
}

func (res Resources) Size() int {
	return res.Rock + res.Wood + res.Food + res.Oil + res.Copper + res.WoodPlank
}

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
	CitizenFeedingRate int        `json:"citizenFeedingRate"`
	EnableEnemies      bool       `json:"enableEnemies"`
	EnableGC           bool       `json:"enableGC"`
	EnemyTickRate      int        `json:"enemyTickRate"`
	FoodGatherRate     int        `json:"foodGatherRate"`
	GcTickRate         int        `json:"gcTickRate"`
	MaxLoad            int        `json:"maxLoad"`
	MaxMissesPerPlayer int        `json:"maxMissesPerPlayer"`
	MaxMissesPerUnit   int        `json:"maxMissesPerUnit"`
	OilGatherRate      int        `json:"oilGatherRate"`
	RockGatherRate     int        `json:"rockGatherRate"`
	ServerId           string     `json:"serverId"`
	TicksPerSeconds    float64    `json:"ticksPerSeconds"`
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
	OpCode    string   `json:"opcode"`
	ReportId  *uid.Uid `json:"reportId"`
	ErrorCode *string  `json:"errorCode"`
	Error     *string  `json:"error"`
	Login     string   `json:"login"`
	UnitId    *uid.Uid `json:"unitId"`
	Misses    int      `json:"misses"`
}
