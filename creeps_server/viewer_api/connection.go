package viewer_api

import (
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

type connection struct {
	isConnected atomic.Bool
	isAdmin     atomic.Bool

	socketLock sync.Mutex
	socket     *websocket.Conn

	chunksLock       sync.RWMutex
	subscribedChunks map[Point]bool

	unitsLock  sync.RWMutex
	knownUnits map[uid.Uid]bool

	playersLock  sync.RWMutex
	knownPlayers map[uid.Uid]bool
}

func (conn *connection) setIsUnitKnown(id uid.Uid, known bool) {
	conn.unitsLock.Lock()
	defer conn.unitsLock.Unlock()
	if known {
		conn.knownUnits[id] = true
	} else {
		delete(conn.knownUnits, id)
	}
}

func (conn *connection) subedToChunk(chunk Point) bool {
	conn.chunksLock.RLock()
	defer conn.chunksLock.RUnlock()
	return conn.subscribedChunks[chunk]
}
