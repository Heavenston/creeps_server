package viewer_api

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
	. "github.com/heavenston/creeps_server/creeps_lib/viewer_api_model"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/server/entities"
	"github.com/rs/zerolog/log"
)

type connection struct {
	viewer *ViewerServer

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

func (conn *connection) sendMessage(kind string, content any) {
	contentbytes, err := json.Marshal(content)
	if err != nil {
		log.Warn().Err(err).Msg("full chunk ser error")
		return
	}
	conn.socketLock.Lock()
	conn.socket.WriteJSON(Message{
		Kind:    kind,
		Content: contentbytes,
	})
	conn.socketLock.Unlock()
}

func (conn *connection) sendChunk(chunk *terrain.Chunk) {
	if !chunk.IsGenerated() {
		return
	}

	tiles := make([]byte, 2*terrain.ChunkSize*terrain.ChunkSize)
	for y := 0; y < terrain.ChunkSize; y++ {
		for x := 0; x < terrain.ChunkSize; x++ {
			i := 2 * (x + y*terrain.ChunkSize)
			tile := chunk.GetTile(Point{X: x, Y: y})
			tiles[i] = byte(tile.Kind)
			tiles[i+1] = byte(tile.Value)
		}
	}
	conn.sendMessage("fullchunk", S2CFullChunk{
		ChunkPos: chunk.GetChunkPos(),
		Tiles:    tiles,
	})
}

func (conn *connection) sendUnit(unit server.IUnit) bool {
	conn.unitsLock.Lock()
	// get the lock for the entire duration to make sure we don't double send
	defer conn.unitsLock.Unlock()

	if conn.knownUnits[unit.GetId()] {
		return false
	}
	conn.sendMessage("unit", S2CUnit{
		OpCode:   unit.GetOpCode(),
		UnitId:   unit.GetId(),
		Owner:    unit.GetOwner(),
		Position: unit.GetPosition(),
		Upgraded: unit.IsUpgraded(),
	})
	conn.knownUnits[unit.GetId()] = true
	return false
}

func (conn *connection) sendPlayer(player *entities.Player) bool {
	conn.playersLock.Lock()
	// get the lock for the entire duration to make sure we don't double send
	defer conn.playersLock.Unlock()

	if conn.knownPlayers[player.GetId()] {
		return false
	}
	conn.sendMessage("playerSpawn", S2CPlayerSpawn{
		Id:            player.GetId(),
		SpawnPosition: player.GetSpawnPoint(),
		Username:      player.GetUsername(),
		Resources:     player.GetResources(),
	})
	conn.knownPlayers[player.GetId()] = true
	return false
}

func (conn *connection) handleServerEvent(event server.IServerEvent) {
	getActionData := func(action *server.Action) ActionData {
		data := ActionData{
			ActionOpCode: action.OpCode,
			ReportId:     action.ReportId,
			Parameter:    action.Parameter,
		}
		return data
	}

	switch e := event.(type) {
	case *server.UnitSpawnEvent:
		conn.sendUnit(e.Unit)
	case *server.UnitDespawnEvent:
		conn.sendMessage("unitDespawned", S2CUnitDespawn{
			UnitId: e.Unit.GetId(),
		})
		conn.setIsUnitKnown(e.Unit.GetId(), false)
	case *server.UnitMovedEvent:
		newChunk := terrain.Global2ContainingChunkCoords(e.To)

		// unit is now out of bounds so we need to make the client forget
		// about it
		if !conn.subedToChunk(newChunk) {
			conn.sendMessage("unitDespawned", S2CUnitDespawn{
				UnitId: e.Unit.GetId(),
			})
			conn.setIsUnitKnown(e.Unit.GetId(), false)
		}
		// the acutal information that the unit moved is done through the
		// report
	case *server.UnitStartedActionEvent:
		// note: we do skip the action...
		if conn.sendUnit(e.Unit) {
			break
		}

		conn.sendMessage("unitStartedAction", S2CUnitStartedAction{
			UnitId: e.Unit.GetId(),
			Action: getActionData(e.Action),
		})
	case *server.UnitFinishedActionEvent:
		// note: if it didn't know about the unit it won't know about
		//       the action
		if conn.sendUnit(e.Unit) {
			break
		}

		content := S2CUnitFinishedAction{
			UnitId: e.Unit.GetId(),
			Action: getActionData(e.Action),
			Report: e.Report,
		}

		conn.sendMessage("unitFinishedAction", content)
	case *entities.PlayerSpawnEvent:
		conn.sendPlayer(e.Player)
	case *entities.PlayerDespawnEvent:
		conn.playersLock.Lock()
		// don't double send
		if !conn.knownPlayers[e.Player.GetId()] {
			conn.playersLock.Unlock()
			break
		}
		conn.sendMessage("playerDespawn", S2CPlayerDespawn{
			Id: e.Player.GetId(),
		})
		delete(conn.knownPlayers, e.Player.GetId())
		conn.playersLock.Unlock()
	}
}
