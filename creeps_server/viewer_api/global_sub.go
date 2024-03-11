package viewer_api

import (
	"encoding/json"
	"time"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	. "github.com/heavenston/creeps_server/creeps_lib/viewer_api_model"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/server/entities"
	"github.com/rs/zerolog/log"
)

func (viewer *ViewerServer) handleClientGlobalEvents(conn *connection) {
	serverEventsChannel := make(chan server.IServerEvent, 2048)
	serverEventsHandle := viewer.Server.Events().Subscribe(serverEventsChannel, AABB{})
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

	sendPlayer := func(player *entities.Player) bool {
		conn.playersLock.Lock()
		// get the lock for the entire duration to make sure we don't double send
		defer conn.playersLock.Unlock()

		if conn.knownPlayers[player.GetId()] {
			return false
		}
		sendMessage("playerSpawn", S2CPlayerSpawn{
			Id:            player.GetId(),
			SpawnPosition: player.GetSpawnPoint(),
			Username:      player.GetUsername(),
			Resources:     player.GetResources(),
		})
		conn.knownPlayers[player.GetId()] = true
		return false
	}

	for _, entity := range viewer.Server.Entities().GetAllIntersects(AABB{}) {
		if player, ok := entity.(*entities.Player); ok {
			sendPlayer(player)
		}
	}

	for {
		if !conn.isConnected.Load() {
			break
		}

		select {
		case event, ok := (<-serverEventsChannel):
			if !ok {
				log.Trace().Msg("server events channel closed")
				break
			}

			if e, ok := event.(*entities.PlayerSpawnEvent); ok {
				sendPlayer(e.Player)
			}
			if e, ok := event.(*entities.PlayerDespawnEvent); ok {
				conn.playersLock.Lock()
				// don't double send
				if !conn.knownPlayers[e.Player.GetId()] {
					conn.playersLock.Unlock()
					break
				}
				sendMessage("playerDespawn", PlayerDespawnContent{
					Id: e.Player.GetId(),
				})
				delete(conn.knownPlayers, e.Player.GetId())
				conn.playersLock.Unlock()
			}
		// makes sure at lease once every 30s we check if we are still subed to
		// the chunk
		case <-time.After(time.Second * 30):
		}
	}
}


