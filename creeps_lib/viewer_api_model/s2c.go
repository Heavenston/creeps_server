package viewerapimodel

import (
	"github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

//tygo:emit
var _ = `
export type S2CMessage =
| { kind: "init"; content: S2CInit; }
| { kind: "fullchunk"; content: S2CFullChunk; }
| { kind: "tileChange"; content: S2CTileChange; }
| { kind: "unit"; content: S2CUnit; }
| { kind: "unitDespawn"; content: S2CUnitDespawn; }
| { kind: "unitStartedAction"; content: S2CUnitStartedAction; }
| { kind: "unitFinishedAction"; content: S2CUnitFinishedAction; }
| { kind: "playerSpawn"; content: S2CPlayerSpawn; }
| { kind: "playerDespawn"; content: S2CPlayerDespawn; }
;
`

type S2CInit struct {
	ChunkSize int                  `json:"chunkSize"`
	Setup     *model.SetupResponse `json:"setup" tstype:",required"`
	Costs     *model.CostsResponse `json:"costs" tstype:",required"`
}

type S2CFullChunk struct {
	ChunkPos geom.Point `json:"chunkPos"`
	// will be base64 encoded
	// each tile has two bytes, one for the kind and one for its value
	// see terrain/tile.go for the correspondance
	// encoded in row-major order
	// can be empty if the chunk isn't generated
	Tiles []byte `json:"tiles"`
}

type S2CTileChange struct {
	TilePos geom.Point `json:"tilePos"`
	Kind    byte       `json:"kind" tstype:"number"`
	Value   byte       `json:"value" tstype:"number"`
}

type S2CUnit struct {
	OpCode   string     `json:"opCode"`
	UnitId   uid.Uid    `json:"unitId"`
	Owner    uid.Uid    `json:"owner"`
	Position geom.Point `json:"position"`
	Upgraded bool       `json:"upgraded"`
}

type S2CUnitDespawn struct {
	UnitId uid.Uid `json:"unitId"`
}

type S2CUnitStartedAction struct {
	UnitId uid.Uid    `json:"unitId"`
	Action ActionData `json:"action"`
}

type S2CUnitFinishedAction struct {
	UnitId uid.Uid       `json:"unitId"`
	Action ActionData    `json:"action"`
	Report model.IReport `json:"report"`
}

type S2CPlayerSpawn struct {
	Id            uid.Uid         `json:"id"`
	SpawnPosition geom.Point      `json:"spawnPosition"`
	Username      string          `json:"username"`
	Resources     model.Resources `json:"resources"`
}

type S2CPlayerDespawn struct {
	Id uid.Uid `json:"id"`
}
