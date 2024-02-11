package model

import "creeps.heav.fr/geom"

type Resources struct {
	Rock int `json:"rock"`
	Wood int `json:"wood"`
	Food int `json:"food"`
	Oil int `json:"oil"`
	Copper int `json:"copper"`
	WoodPlanks int `json:"woodPlanks"`
}

type StatusResponse struct {
	Running bool `json:"running"`
}

type StatisticsResponse struct {
	ServerId string `json:"qerverId"`
	GameRunning string `json:"qameRunning"`
	Tick int `json:"qick"`
	Dimensions geom.Point `json:"qimensions"`
	Players []struct {
		Name string `json:"name"`
		Status string `json:"status"`
		Units int `json:"units"`
		Buildings int `json:"buildings"`
		Resources Resources  `json:"resources"`
		Achievements []string `json:"achievements"`
	} `json:"players"`
}

type CostResponse struct  {
	Resources
	Cast int
}

type CostsResponse struct {
	BuildHousehold CostResponse `json:"buildHousehold"`
	BuildRoad CostResponse `json:"buildRoad"`
	BuildSawmill CostResponse `json:"buildSawmill"`
	BuildSmeltery CostResponse `json:"buildSmeltery"`
	BuildTownHall CostResponse `json:"buildTownHall"`
	Dismantle CostResponse `json:"dismantle"`
	Farm CostResponse `json:"farm"`
	FetchMessage CostResponse `json:"fetchMessage"`
	FireBomberBot CostResponse `json:"fireBomberBot"`
	FireTurret CostResponse `json:"fireTurret"`
	Gather CostResponse `json:"gather"`
	Move CostResponse `json:"move"`
	Noop CostResponse `json:"noop"`
	Observe CostResponse `json:"observe"`
	RefineCopper CostResponse `json:"refineCopper"`
	RefineWoodPlank CostResponse `json:"refineWoodPlank"`
	SendMessage CostResponse `json:"sendMessage"`
	SpawnBomberBot CostResponse `json:"spawnBomberBot"`
	SpawnTurret CostResponse `json:"spawnTurret"`
	Unload CostResponse `json:"unload"`
	UpgradeBomberBot CostResponse `json:"upgradeBomberBot"`
	UpgradeCitizen CostResponse `json:"upgradeCitizen"`
	UpgradeTurret CostResponse `json:"upgradeTurret"`
}

type SetupResponse struct {
	CitizenFeedingRate int `json:"citizenFeedingRate"`
	EnableEnemies bool `json:"enableEnemies"`
	EnableGC bool `json:"enableGC"`
	EnemyTickRate int `json:"enemyTickRate"`
	FoodGatherRate int `json:"foodGatherRate"`
	GcTickRate int `json:"gcTickRate"`
	MaxLoad int `json:"maxLoad"`
	MaxMissesPerPlayer int `json:"maxMissesPerPlayer"`
	MaxMissesPerUnit int `json:"maxMissesPerUnit"`
	OilGatherRate int `json:"oilGatherRate"`
	RockGatherRate int `json:"rockGatherRate"`
	ServerId string `json:"serverId"`
	TicksPerSeconds float64 `json:"ticksPerSeconds"`
	TrackAchievements bool `json:"trackAchievements"`
	WoodGatherRate int `json:"woodGatherRate"`
	WorldDimension geom.Point `json:"worldDimension"`
}

type InitResponse struct {
	Citizen1Id *string `json:"citizen1Id"`
	Citizen2Id *string `json:"citizen2Id"`
	Costs *CostsResponse `json:"costs"`
	Error *string `json:"error"`
	HouseholdCoordinates *geom.Point `json:"householdCoordinates"`
	Login string `json:"login"`
	PlayerId *int16 `json:"playerId"`
	Resources *Resources `json:"resources"`
	Setup *SetupResponse `json:"setup"`
	Tick *int `json:"tick"`
	TownHallCoordinates *geom.Point `json:"townHallCoordinates"`
}

type CommandResponse struct {
	OpCode *string `json:"opCode"`
	ReportId *string `json:"reportId"`
	ErrorCode *string `json:"errorCode"`
	Error *string `json:"error"`
	Login *string `json:"login"`
	UnitId *string `json:"unitId"`
	Misses int `json:"misses"`
}

type Report struct {
	OpCode *string `json:"opCode"`
	ReportId *string `json:"reportId"`
	UnitId *string `json:"unitId"`
	Login *string `json:"login"`
	UnitPosition *geom.Point `json:"unitPosition"`
	Status *string `json:"status"`
	ErrorCode *string `json:"errorCode"`
	Error *string `json:"error"`
}
