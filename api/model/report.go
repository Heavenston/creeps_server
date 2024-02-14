package model

import (
	"reflect"

	"creeps.heav.fr/geom"
)

// used by some reports
type Unit struct {
	OpCode   string     `json:"opcode"`
	OwnerId  string     `json:"player"`
	Position geom.Point `json:"position"`
}

// used by some reports
type Building struct {
	OpCode   string     `json:"opcode"`
	OwnerId  string     `json:"ownerId"`
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
	OpCode       *string     `json:"opcode"`
	ReportId     *string     `json:"reportId"`
	UnitId       *string     `json:"unitId"`
	Login        *string     `json:"login"`
	UnitPosition *geom.Point `json:"unitPosition"`
	Status       *string     `json:"status"`
	ErrorCode    *string     `json:"errorCode"`
	Error        *string     `json:"error"`
}

func (r *Report) GetReport() *Report {
	return r
}

func (r *Report) GetParameterType() reflect.Type {
	return nil
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
	Resources     string `json:"resources"`
	Gathered      int    `json:"gathered"`
	ResourcesLeft int    `json:"resourcesLeft"`
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
	SpawnedCitizen1Id string `json:"spawnedCitizen1Id"`
	SpawnedCitizen2Id string `json:"spawnedCitizen2Id"`
}

type SpawnReport struct {
	Report
	SpawnedUnitId string `json:"spawnedUnitId"`
	SpawnedUnit   Unit   `json:"spawnedUnit"`
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
	Target geom.Point `json:"target"`
	KilledUnits []Unit `json:"killedUnits"`
}

func (r *FireReport) GetParameterType() reflect.Type {
	return reflect.TypeFor[FireParameter]()
}
