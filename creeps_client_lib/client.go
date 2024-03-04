package creepsclientlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
)

type lockedResources struct {
	lock      sync.RWMutex
	resources model.Resources
}

type Client struct {
	serverAddr string
	apiPrefix  string
	login      string

	tilemap      atomic.Pointer[terrain.Tilemap]
	initResponse atomic.Pointer[model.InitResponse]

	playerResources    *model.AtomicResources
	unitsResources     sync.Map

	unitsPositions     sync.Map
}

// error returned by Get*Report methods if they get a model.ReportError response
type ReportError struct {
	Report *model.ErrorReport
}

func (err *ReportError) Error() string {
	return err.Report.ErrorCode
}

type ErrCommand struct {
	response *model.CommandResponse
}

func (err *ErrCommand) Error() string {
	return *err.response.Error
}

type ErrNotEnoughResources struct {
}

func (err *ErrNotEnoughResources) Error() string {
	return "Not enough resources :("
}

// Creates a new client, makes no requests
//
// example:
// client := NewClient("localhost:1664", "heavenstone")
func NewClient(serverAddr string, login string) *Client {
	client := new(Client)

	client.serverAddr = serverAddr
	client.apiPrefix = "http://" + serverAddr
	client.login = login
	client.playerResources = &model.AtomicResources{}

	return client
}

func (client *Client) ServerAddr() string {
	return client.serverAddr
}

func (client *Client) ApiPrefix() string {
	return client.apiPrefix
}

func (client *Client) Login() string {
	return client.login
}

func (client *Client) InitResponse() *model.InitResponse {
	return client.initResponse.Load()
}

// Computes the tick duration from the init response
func (client *Client) TickDuration() time.Duration {
	return time.Second / time.Duration(client.initResponse.Load().Setup.TicksPerSecond)
}

func (client *Client) SleepFor(ticks int) {
	time.Sleep(client.TickDuration()*time.Duration(ticks) + time.Millisecond*5)
}

func (client *Client) SetTilemap(tm *terrain.Tilemap) {
	client.tilemap.Store(tm)
}

func (client *Client) PlayerResources() *model.AtomicResources {
	return client.playerResources
}

func (client *Client) UnitResources(unitId uid.Uid) *model.AtomicResources {
	val, _ := client.unitsResources.LoadOrStore(unitId, &model.AtomicResources{})
	return val.(*model.AtomicResources)
}

func (client *Client) UnitPosition(unitId uid.Uid) *AtomicPoint {
	val, loaded := client.unitsPositions.LoadOrStore(unitId, &AtomicPoint{})
	pos := val.(*AtomicPoint)
	if !loaded {
		if unitId == *client.InitResponse().Citizen1Id {
			pos.Store(*client.InitResponse().HouseholdCoordinates)
		}
		if unitId == *client.InitResponse().Citizen2Id {
			pos.Store(*client.InitResponse().HouseholdCoordinates)
		}
	}
	return pos
}

func (client *Client) RawGet(url string) (*http.Response, error) {
	return http.Get(client.apiPrefix + url)
}

func (client *Client) Get(url string, responseDest any) error {
	resp, err := client.RawGet(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseDest)
	if err != nil {
		return err
	}

	return nil
}

// makes a post request by json encoding the given body and parsing the response
// into responseDest
// reqBody can be nil
func (client *Client) Post(url string, responseDest any, reqBody any) error {
	var reqBodyReader io.Reader
	if reqBody != nil {
		reqBodyEncoded, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		reqBodyReader = bytes.NewReader(reqBodyEncoded)
	}

	resp, err := http.Post(client.apiPrefix+url, "application/json", reqBodyReader)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseDest)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) GetStatus() (resp model.StatusResponse, err error) {
	err = client.Get("/status", &resp)
	return
}

func (client *Client) GetStatistics() (resp model.StatisticsResponse, err error) {
	err = client.Get("/statistics", &resp)
	return
}

// The response is also stored inside the client
func (client *Client) PostInit() (resp *model.InitResponse, err error) {
	err = client.Post("/init/"+client.login, &resp, nil)
	client.initResponse.Store(resp)
	if resp != nil && resp.Resources != nil {
		client.PlayerResources().Store(*resp.Resources)
	}
	return
}

// Tried to post a command that requires a body without giving a body
type ErrCommandRequiresBody struct {
	OpCode model.ActionOpCode
}

func (err *ErrCommandRequiresBody) Error() string {
	return fmt.Sprintf(
		"OpCode %s requires a body, please use PostCommandWithBody",
		err.OpCode,
	)
}

// If the opcode requires a parameter this will return a ErrCommandRequiresBody
// so please use PostCommandWithBody if applicable
func (client *Client) PostCommand(
	unitId uid.Uid,
	opcode model.ActionOpCode,
) (resp *model.CommandResponse, err error) {
	if opcode.ParameterType() != nil {
		err = &ErrCommandRequiresBody{
			OpCode: opcode,
		}
		return
	}
	err = client.Post(
		"/command/"+client.login+"/"+string(unitId)+"/"+string(opcode),
		&resp,
		nil,
	)
	if resp != nil && resp.ErrorCode != nil {
		err = &ErrCommand{
			response: resp,
		}
		resp = nil
	}
	return
}

// Like PostCommand but adds the given body serialized in json
func (client *Client) PostCommandWithBody(
	unitId uid.Uid,
	opcode model.ActionOpCode,
	body any,
) (resp *model.CommandResponse, err error) {
	err = client.Post(
		"/command/"+client.login+"/"+string(unitId)+"/"+string(opcode),
		&resp,
		body,
	)
	return
}

// Gets the report and fills the given variable
// If the server responds with an error report a ReportError is returned
//
// expample:
// var report model.SpawnReport
// err := client.GetReport(id, &report)
//
//	if err != nil {
//	    return err
//	}
func (client *Client) GetReport(
	reportId uid.Uid,
	reportOut any,
) error {
	resp, err := client.RawGet("/report/" + string(reportId))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errResp model.ErrorReport
	err = json.Unmarshal(body, &errResp)
	if err == nil && errResp.Status == "ERROR" {
		return &ReportError{Report: &errResp}
	}

	err = json.Unmarshal(body, reportOut)
	if err != nil {
		return err
	}

	if rep, ok := reportOut.(model.IReport); ok {
		RegisterReport(client, rep)
	}

	return nil
}

func (client *Client) Command(
	unitId uid.Uid,
	opcode model.ActionOpCode,
	upgradeCost *model.CostResponse,
) (model.IReport, error) {
	cost := opcode.GetCost(client.InitResponse().Costs, upgradeCost)

	if !client.playerResources.TrySub(cost.Resources) {
		return nil, &ErrNotEnoughResources{}
	}

	cmd, err := client.PostCommand(unitId, opcode)
	if err != nil {
		// post fail so we credit the resources back
		client.playerResources.Add(cost.Resources)
		return nil, err
	}

	client.SleepFor(cost.Cast)

	report := reflect.New(opcode.GetReportType())
	err = client.GetReport(*cmd.ReportId, report.Interface())
	if err != nil {
		return nil, err
	}
	return report.Interface().(model.IReport), nil
}
