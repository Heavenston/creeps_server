package viewer

import (
	"encoding/json"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/uid"
)

type message struct {
	Kind    string          `json:"kind"`
	Content json.RawMessage `json:"content"`
}

// first packet sent by the server to the client as soon as the connection is
// established with various informations
type initContent struct {
	ChunkSize int                  `json:"chunkSize"`
	Setup     *model.SetupResponse `json:"setup"`
	Costs     *model.CostsResponse `json:"costs"`
}

type fullChunkContent struct {
	ChunkPos Point `json:"chunkPos"`
	// will be base64 encoded
	// each tile has two bytes, one for the kind and one for its value
	// see terrain/tile.go for the correspondance
	// encoded in row-major order
	// can be empty if the chunk isn't generated
	Tiles []byte `json:"tiles"`
}

type tileChangeContent struct {
	TilePos Point `json:"tilePos"`
	Kind    byte  `json:"kind"`
	Value   byte  `json:"value"`
}

// sent by the server when a unit spawned
type unitContent struct {
	OpCode   string  `json:"opCode"`
	UnitId   uid.Uid `json:"unitId"`
	Owner    uid.Uid `json:"owner"`
	Position Point   `json:"position"`
	Upgraded bool    `json:"upgraded"`
}

type unitMovementContent struct {
	UnitId uid.Uid `json:"unitId"`
	New    Point   `json:"new"`
}

type unitUpgradedContent struct {
	UnitId uid.Uid `json:"unitId"`
}

// sent by the server when a unit dies or gets out of the subscribed chunks
type unitDespawnContent struct {
	UnitId uid.Uid `json:"unitId"`
}

type playerSpawnContent struct {
	Id            uid.Uid         `json:"id"`
	SpawnPosition Point           `json:"spawnPosition"`
	Username      string          `json:"username"`
	Resources     model.Resources `json:"resources"`
}

type playerDespawnContent struct {
	Id uid.Uid `json:"id"`
}

// sent by the front end to subscribe to a chunk content
type subscribeRequestContent struct {
	ChunkPos Point `json:"chunkPos"`
}

// sent by the front end to unsubscribe from a chunk content
type unsubscribeRequestContent struct {
	ChunkPos Point `json:"chunkPos"`
}
