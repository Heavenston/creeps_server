package viewerapimodel

import (
	"github.com/heavenston/creeps_server/creeps_lib/geom"
)

//tygo:emit
var _ = `
export type C2SMessage =
| { kind: "init"; content: C2SInit; }
| { kind: "subscribe"; content: C2SSubscribeRequest; }
| { kind: "unsubscribe"; content: C2SUnsubscribeRequest; }
;
`

type C2SInit struct {
	AuthPassword string `json:"auth_password,omitempty"`
}

// Subscribs to all entities and tile changes from the given chunk
// after receiving, if available, the server sends the full chunk tiles
type C2SSubscribeRequest struct {
	ChunkPos geom.Point `json:"chunkPos"`
}

// cancels a previous subscribe request
type C2SUnsubscribeRequest struct {
	ChunkPos geom.Point `json:"chunkPos"`
}
