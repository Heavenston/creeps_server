package viewer_api

import (
	"time"

	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/server/entities"
	"github.com/rs/zerolog/log"
)

func (conn *connection) handleClientGlobalEvents() {
	viewer := conn.viewer

	serverEventsChannel := make(chan server.IServerEvent, 2048)
	serverEventsHandle := viewer.Server.Events().
		Subscribe(serverEventsChannel, spatialmap.Extent{ IsGlobal: true })
	defer serverEventsHandle.Cancel()

	for _, entity := range viewer.Server.Entities().
		GetAllCollides(spatialmap.Extent{ IsGlobal: true }) {
		if player, ok := entity.(*entities.Player); ok {
			conn.sendPlayer(player)
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

			conn.handleServerEvent(event)
		// makes sure at lease once every 30s we check if we are still subed to
		// the chunk
		case <-time.After(time.Second * 30):
		}
	}
}
