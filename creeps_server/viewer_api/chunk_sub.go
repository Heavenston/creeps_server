package viewer_api

import (
	"time"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
	. "github.com/heavenston/creeps_server/creeps_lib/viewer_api_model"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/rs/zerolog/log"
)

func (conn *connection) handleClientSubscription(chunkPos Point) {
	viewer := conn.viewer

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

	conn.sendChunk(chunk)

	for _, entity := range viewer.Server.Entities().GetAllIntersects(aabb) {
		if unit, ok := entity.(server.IUnit); ok {
			conn.sendUnit(unit)
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
				conn.sendMessage("unitDespawned", S2CUnitDespawn{
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
				conn.sendMessage("tileChange", S2CTileChange{
					TilePos: change.UpdatedPosition.Add(chunkPos.Times(terrain.ChunkSize)),
					Kind:    byte(change.NewValue.Kind),
					Value:   change.NewValue.Value,
				})
			}
			if _, ok := event.(terrain.GeneratedChunkEvent); ok {
				conn.sendChunk(chunk)
			}
		case event, ok := (<-serverEventsChannel):
			if !ok {
				log.Trace().Msg("server events channel closed")
				break
			}

			conn.handleServerEvent(event)
		// makes sure at lease once every 30s we check if we are still subed to
		// the chunk
		case <-time.After(time.Second * 30):
		}
	}
}
