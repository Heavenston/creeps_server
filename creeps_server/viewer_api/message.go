package viewer_api

// c2s is client to server
// s2c is server to client

import (
	"encoding/json"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

type message struct {
	Kind    string          `json:"kind"`
	Content json.RawMessage `json:"content"`
}

type c2sInit struct {
	
}

// Sent as a response to c2sInit
type s2cInit struct {
	ChunkSize int                  `json:"chunkSize"`
	Setup     *model.SetupResponse `json:"setup"`
	Costs     *model.CostsResponse `json:"costs"`
}

type s2cFullChunk struct {
	ChunkPos Point `json:"chunkPos"`
	// will be base64 encoded
	// each tile has two bytes, one for the kind and one for its value
	// see terrain/tile.go for the correspondance
	// encoded in row-major order
	// can be empty if the chunk isn't generated
	Tiles []byte `json:"tiles"`
}

type s2cTileChange struct {
	TilePos Point `json:"tilePos"`
	Kind    byte  `json:"kind"`
	Value   byte  `json:"value"`
}

type actionData struct {
	ActionOpCode model.ActionOpCode `json:"actionOpCode"`
	ReportId     uid.Uid            `json:"reportId"`
	Parameter    any                `json:"parameter,omitempty"`
}

type s2cUnit struct {
	OpCode   string  `json:"opCode"`
	UnitId   uid.Uid `json:"unitId"`
	Owner    uid.Uid `json:"owner"`
	Position Point   `json:"position"`
	Upgraded bool    `json:"upgraded"`
}

// sent by the server when a unit dies or gets out of the subscribed chunks
type s2cUnitDespawn struct {
	UnitId uid.Uid `json:"unitId"`
}

type s2cUnitStartedAction struct {
	UnitId uid.Uid    `json:"unitId"`
	Action actionData `json:"action"`
}

type s2cUnitFinishedAction struct {
	UnitId uid.Uid       `json:"unitId"`
	Action actionData    `json:"action"`
	Report model.IReport `json:"report"`
}

type s2cPlayerSpawn struct {
	Id            uid.Uid         `json:"id"`
	SpawnPosition Point           `json:"spawnPosition"`
	Username      string          `json:"username"`
	Resources     model.Resources `json:"resources"`
}

type playerDespawnContent struct {
	Id uid.Uid `json:"id"`
}

// Subscribs to all entities and tile changes from the given chunk
// after receiving, if available, the server sends the full chunk tiles
type c2sSubscribeRequest struct {
	ChunkPos Point `json:"chunkPos"`
}

// cancels a previous subscribe request
type c2sUnsubscribeRequest struct {
	ChunkPos Point `json:"chunkPos"`
}
