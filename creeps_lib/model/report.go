package model

import (
	"github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

// used by some reports
type Unit struct {
	OpCode   string     `json:"opcode"`
	Player   string     `json:"player"`
	Position geom.Point `json:"position"`
}

// used by some reports
type Building struct {
	OpCode   string     `json:"opcode"`
	Player   uint16     `json:"ownerId"`
	Position geom.Point `json:"position"`
}

// used by message fetch
type Message struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

//tygo:ignore
type IReport interface {
	GetReport() *Report
}

type Report struct {
	OpCode       ActionOpCode `json:"opcode"`
	ReportId     uid.Uid      `json:"reportId"`
	UnitId       uid.Uid      `json:"unitId"`
	Login        string       `json:"login"`
	UnitPosition geom.Point   `json:"unitPosition"`
	Status       string       `json:"status"`
}

func (r *Report) GetReport() *Report {
	return r
}

type ErrorReport struct {
	Report    `json:",inline" tstype:",extends"`
	ErrorCode string `json:"errorCode"`
	Error     string `json:"-"`
}

type NoOpReport struct {
	Report `json:",inline" tstype:",extends"`
}

type ObserveReport struct {
	Report `json:",inline" tstype:",extends"`
	Tiles  []uint16 `json:"tiles"`
	Units  []Unit   `json:"units"`
}

type MoveReport struct {
	ObserveReport `json:",inline" tstype:",extends"`
	NewPosition   geom.Point `json:"newPosition"`
}

type GatherReport struct {
	Report        `json:",inline" tstype:",extends"`
	Resource      ResourceKind `json:"resource"`
	Gathered      int          `json:"gathered"`
	ResourcesLeft int          `json:"resourcesLeft"`
}

type UnloadReport struct {
	Report            `json:",inline" tstype:",extends"`
	CreditedResources Resources `json:"creditedResources"`
}

type FarmReport struct {
	Report       `json:",inline" tstype:",extends"`
	FoodQuantity int `json:"foodQuantity"`
}

type BuildReport struct {
	Report   `json:",inline" tstype:",extends"`
	Building Building `json:"building"`
}

type BuildHouseHoldReport struct {
	BuildReport       `json:",inline" tstype:",extends"`
	SpawnedCitizen1Id uid.Uid `json:"spawnedCitizen1Id"`
	SpawnedCitizen2Id uid.Uid `json:"spawnedCitizen2Id"`
}

type SpawnReport struct {
	Report        `json:",inline" tstype:",extends"`
	SpawnedUnitId uid.Uid `json:"spawnedUnitId"`
	SpawnedUnit   Unit    `json:"spawnedUnit"`
}

type DismantleReport struct {
	Report `json:",inline" tstype:",extends"`
}

type UpgradeReport struct {
	Report `json:",inline" tstype:",extends"`
}

type RefineReport struct {
	Report `json:",inline" tstype:",extends"`
}

type MessageSendReport struct {
	Report    `json:",inline" tstype:",extends"`
	Recipient string `json:"recipient"`
}

type MessageFetchReport struct {
	Report          `json:",inline" tstype:",extends"`
	FetchedMessages []Message `json:"fetchedMessages"`
}

type FireReport struct {
	Report      `json:",inline" tstype:",extends"`
	Target      geom.Point `json:"target"`
	KilledUnits []Unit     `json:"killedUnits"`
}
