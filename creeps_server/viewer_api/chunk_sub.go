package viewer_api

import (
	"encoding/json"
	"time"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
	. "github.com/heavenston/creeps_server/creeps_lib/viewer_api_model"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/rs/zerolog/log"
)

func (viewer *ViewerServer) handleClientSubscription(
	chunkPos Point,
	conn *connection,
) {
	chunk := viewer.Server.Tilemap().CreateChunk(chunkPos)

	terrainChangeChannel := make(chan any, 2048)
	terrainCancelHandle := chunk.UpdatedEventProvider.Subscribe(terrainChangeChannel)
	defer terrainCancelHandle.Cancel()

	serverEventsChannel := make(chan server.IServerEvent, 2048)
	aabb := AABB{
		From: chunkPos.Times(terrain.ChunkSize),
		Size: Point{
			// +1 seems to fix some missed events
			X: terrain.ChunkSize + 1,
			Y: terrain.ChunkSize + 1,
		},
	}
	serverEventsHandle := viewer.Server.Events().Subscribe(serverEventsChannel, aabb)
	defer serverEventsHandle.Cancel()

	sendMessage := func(kind string, content any) {
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

	// send full chunk
	sendTerrain := func() {
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
		sendMessage("fullchunk", S2CFullChunk{
			ChunkPos: chunkPos,
			Tiles:    tiles,
		})
	}

	sendUnit := func(unit server.IUnit) bool {
		conn.unitsLock.Lock()
		// get the lock for the entire duration to make sure we don't double send
		defer conn.unitsLock.Unlock()

		if conn.knownUnits[unit.GetId()] {
			return false
		}
		sendMessage("unit", S2CUnit{
			OpCode:   unit.GetOpCode(),
			UnitId:   unit.GetId(),
			Owner:    unit.GetOwner(),
			Position: unit.GetPosition(),
			Upgraded: unit.IsUpgraded(),
		})
		conn.knownUnits[unit.GetId()] = true
		return false
	}

	getActionData := func(action *server.Action) ActionData {
		data := ActionData{
			ActionOpCode: action.OpCode,
			ReportId:     action.ReportId,
			Parameter:    action.Parameter,
		}
		return data
	}

	sendTerrain()

	for _, entity := range viewer.Server.Entities().GetAllIntersects(aabb) {
		if unit, ok := entity.(server.IUnit); ok {
			sendUnit(unit)
		}
	}

	for {
		if !conn.isConnected.Load() {
			break
		}
		stillSubed := conn.subedToChunk(chunkPos)
		if !stillSubed {
			conn.unitsLock.Lock()

			for _, entity := range viewer.Server.Entities().GetAllIntersects(aabb) {
				id := entity.GetId()
				if !conn.knownUnits[id] {
					continue
				}
				sendMessage("unitDespawned", S2CUnitDespawn{
					UnitId: id,
				})
				delete(conn.knownUnits, id)
			}

			conn.unitsLock.Unlock()
			break
		}

		select {
		case event, ok := (<-terrainChangeChannel):
			if !ok {
				log.Trace().Msg("terrain channel closed")
				break
			}

			if change, ok := event.(terrain.TileUpdateChunkEvent); ok {
				sendMessage("tileChange", S2CTileChange{
					TilePos: change.UpdatedPosition.Add(chunkPos.Times(terrain.ChunkSize)),
					Kind:    byte(change.NewValue.Kind),
					Value:   change.NewValue.Value,
				})
			}
			if _, ok := event.(terrain.GeneratedChunkEvent); ok {
				sendTerrain()
			}
		case event, ok := (<-serverEventsChannel):
			if !ok {
				log.Trace().Msg("server events channel closed")
				break
			}

			if e, ok := event.(*server.UnitSpawnEvent); ok {
				sendUnit(e.Unit)
			}
			if e, ok := event.(*server.UnitDespawnEvent); ok {
				sendMessage("unitDespawned", S2CUnitDespawn{
					UnitId: e.Unit.GetId(),
				})
				conn.setIsUnitKnown(e.Unit.GetId(), false)
			}
			if e, ok := event.(*server.UnitMovedEvent); ok {
				newChunk := terrain.Global2ContainingChunkCoords(e.To)
				if newChunk != chunkPos && !conn.subedToChunk(newChunk) {
					sendMessage("unitDespawned", S2CUnitDespawn{
						UnitId: e.Unit.GetId(),
					})
					conn.setIsUnitKnown(e.Unit.GetId(), false)
					break
				}
			}
			if e, ok := event.(*server.UnitStartedActionEvent); ok {
				// note: we do skip the action...
				if sendUnit(e.Unit) {
					break
				}

				sendMessage("unitStartedAction", S2CUnitStartedAction{
					UnitId: e.Unit.GetId(),
					Action: getActionData(e.Action),
				})
			}
			if e, ok := event.(*server.UnitFinishedActionEvent); ok {
				// note: if it didn't know about the unit it won't know about
				//       the action
				if sendUnit(e.Unit) {
					break
				}

				content := S2CUnitFinishedAction{
					UnitId: e.Unit.GetId(),
					Action: getActionData(e.Action),
					Report: e.Report,
				}

				sendMessage("unitFinishedAction", content)
			}
		// makes sure at lease once every 30s we check if we are still subed to
		// the chunk
		case <-time.After(time.Second * 30):
		}
	}
}

