package model

import (
	"reflect"

	"creeps.heav.fr/geom"
	"creeps.heav.fr/uid"
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

type IReport interface {
	GetReport() *Report
	GetParameterType() reflect.Type
}

type Report struct {
	OpCode       string     `json:"opcode"`
	ReportId     uid.Uid    `json:"reportId"`
	UnitId       uid.Uid    `json:"unitId"`
	Login        string     `json:"login"`
	UnitPosition geom.Point `json:"unitPosition"`
	Status       string     `json:"status"`
}

func (r *Report) GetReport() *Report {
	return r
}

func (r *Report) GetParameterType() reflect.Type {
	return nil
}

type ErrorReport struct {
	Report
	ErrorCode string `json:"errorCode"`
	Error     string `json:"-"`
}

type NoOpReport struct {
	Report
}

type ObserveReport struct {
	Report
	Tiles []uint16 `json:"tiles"`
	Units []Unit   `json:"units"`
}

type MoveReport struct {
	ObserveReport
	NewPosition geom.Point `json:"newPosition"`
}

type GatherReport struct {
	Report
	Resource      ResourceKind `json:"resource"`
	Gathered      int          `json:"gathered"`
	ResourcesLeft int          `json:"resourcesLeft"`
}

type UnloadReport struct {
	Report
	CreditedResources Resources `json:"creditedResources"`
}

type FarmReport struct {
	Report
	FoodQuantity int `json:"foodQuantity"`
}

type BuildReport struct {
	Report
	Building Building `json:"building"`
}

type BuildHouseHoldReport struct {
	BuildReport
	SpawnedCitizen1Id uid.Uid `json:"spawnedCitizen1Id"`
	SpawnedCitizen2Id uid.Uid `json:"spawnedCitizen2Id"`
}

type SpawnReport struct {
	Report
	SpawnedUnitId uid.Uid `json:"spawnedUnitId"`
	SpawnedUnit   Unit    `json:"spawnedUnit"`
}

type DismantleReport struct {
	Report
}

type UpgradeReport struct {
	Report
}

type RefineReport struct {
	Report
}

type MessageSendReport struct {
	Report
	Recipient string `json:"recipient"`
}

func (r *MessageSendReport) GetParameterType() reflect.Type {
	return reflect.TypeFor[MessageSendParameter]()
}

type MessageFetchReport struct {
	Report
	FetchedMessages []Message `json:"fetchedMessages"`
}

type FireReport struct {
	Report
	Target      geom.Point `json:"target"`
	KilledUnits []Unit     `json:"killedUnits"`
}

func (r *FireReport) GetParameterType() reflect.Type {
	return reflect.TypeFor[FireParameter]()
}
