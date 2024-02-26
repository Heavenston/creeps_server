package viewer

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/server/entities"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ViewerServer struct {
	Server *server.Server
	Addr   string
}

type connection struct {
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

func (viewer *ViewerServer) handleClientSubscription(
	chunkPos Point,
	conn *connection,
) {
	chunk := viewer.Server.Tilemap().CreateChunk(chunkPos)

	terrainChangeChannel := make(chan any, 64)
	terrainCancelHandle := chunk.UpdatedEventProvider.Subscribe(terrainChangeChannel)
	defer terrainCancelHandle.Cancel()

	serverEventsChannel := make(chan server.IServerEvent, 64)
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
		conn.socket.WriteJSON(message{
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
		sendMessage("fullchunk", fullChunkContent{
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
		sendMessage("unit", unitContent{
			OpCode:   unit.GetOpCode(),
			UnitId:   unit.GetId(),
			Owner:    unit.GetOwner(),
			Position: unit.GetPosition(),
			Upgraded: unit.IsUpgraded(),
		})
		conn.knownUnits[unit.GetId()] = true
		return false
	}

	sendPlayer := func(player *entities.Player) bool {
		conn.playersLock.Lock()
		// get the lock for the entire duration to make sure we don't double send
		defer conn.playersLock.Unlock()

		if conn.knownPlayers[player.GetId()] {
			return false
		}
		sendMessage("playerSpawn", playerSpawnContent{
			Id:            player.GetId(),
			SpawnPosition: player.GetSpawnPoint(),
			Username:      player.GetUsername(),
			Resources:     player.GetResources(),
		})
		conn.knownPlayers[player.GetId()] = true
		return false
	}

	getActionData := func (action *server.Action) actionData {
		data := actionData{
			ActionOpCode: action.OpCode,
			ReportId: action.ReportId,
			Parameter: action.Parameter,
		}
		return data
	}

	sendTerrain()

	for _, entity := range viewer.Server.Entities().GetAllIntersects(aabb) {
		if unit, ok := entity.(server.IUnit); ok {
			sendUnit(unit)
		}
		if player, ok := entity.(*entities.Player); ok {
			sendPlayer(player)
		}
	}

	for {
		stillSubed := conn.subedToChunk(chunkPos)
		if !stillSubed {
			conn.unitsLock.Lock()

			for _, entity := range viewer.Server.Entities().GetAllIntersects(aabb) {
				id := entity.GetId()
				if !conn.knownUnits[id] {
					continue
				}
				sendMessage("unitDespawned", unitDespawnContent{
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
				sendMessage("tileChange", tileChangeContent{
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
				sendMessage("unitDespawned", unitDespawnContent{
					UnitId: e.Unit.GetId(),
				})
				conn.setIsUnitKnown(e.Unit.GetId(), false)
			}
			if e, ok := event.(*server.UnitMovedEvent); ok {
				newChunk := terrain.Global2ContainingChunkCoords(e.To)
				if newChunk != chunkPos && !conn.subedToChunk(newChunk) {
					sendMessage("unitDespawned", unitDespawnContent{
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

				sendMessage("unitStartedAction", unitStartedActionContent{
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

				content := unitFinishedActionContent{
					UnitId: e.Unit.GetId(),
					Action: getActionData(e.Action),
					Report: e.Report,
				}

				sendMessage("unitFinishedAction", content)
			}
			if e, ok := event.(*entities.PlayerSpawnEvent); ok {
				sendPlayer(e.Player)
			}
			if e, ok := event.(*entities.PlayerDespawnEvent); ok {
				conn.playersLock.Lock()
				// don't double send
				if !conn.knownPlayers[e.Player.GetId()] {
					break
				}
				sendMessage("playerDespawn", playerDespawnContent{
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

func (viewer *ViewerServer) handleClient(conn *websocket.Conn) {
	var err error = nil

	defer conn.Close()
	defer log.Debug().Any("addr", conn.RemoteAddr()).Msg("Websocket connection closed")

	log.Debug().Any("addr", conn.RemoteAddr()).Msg("New websocket connection")

	connection := connection{
		socket:           conn,
		subscribedChunks: make(map[Point]bool),
		knownUnits:       make(map[uid.Uid]bool),
		knownPlayers:     make(map[uid.Uid]bool),
	}

	{
		var initMessage message
		initMessage.Kind = "init"
		var messContent initContent
		messContent.ChunkSize = terrain.ChunkSize
		messContent.Costs = viewer.Server.GetCosts()
		messContent.Setup = viewer.Server.GetSetup()

		bytes, err := json.Marshal(messContent)
		if err != nil {
			goto error
		}
		initMessage.Content = bytes

		bytes, err = json.Marshal(initMessage)
		if err != nil {
			goto error
		}

		err = conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			goto error
		}
	}

	for {
		var mess message
		//t, reader, err := conn.NextReader()
		err = conn.ReadJSON(&mess)
		if err != nil {
			goto error
		}

		log.Trace().
			Any("message_kind", mess.Kind).
			Any("addr", conn.RemoteAddr()).
			Msg("Websocket received message")

		switch mess.Kind {
		case "subscribe":
			var content subscribeRequestContent
			err = json.Unmarshal(mess.Content, &content)
			if err != nil {
				goto softerror
			}

			connection.chunksLock.Lock()

			connection.subscribedChunks[content.ChunkPos] = true

			log.Trace().Any("addr", conn.RemoteAddr()).
				Any("chunkPos", content.ChunkPos).
				Int("subCount", len(connection.subscribedChunks)).
				Msg("Subscribed to a chunk")

			connection.chunksLock.Unlock()
			go viewer.handleClientSubscription(content.ChunkPos, &connection)
		case "unsubscribe":
			var content unsubscribeRequestContent
			err = json.Unmarshal(mess.Content, &content)
			if err != nil {
				goto softerror
			}

			connection.chunksLock.Lock()
			delete(connection.subscribedChunks, content.ChunkPos)

			log.Trace().Any("addr", conn.RemoteAddr()).
				Any("chunkPos", content.ChunkPos).
				Int("subCount", len(connection.subscribedChunks)).
				Msg("Unsubscribed to a chunk")

			connection.chunksLock.Unlock()
		default:
			log.Debug().
				Any("addr", conn.RemoteAddr()).
				Any("mess", mess).
				Any("kind", mess.Kind).
				Msg("Unknown message")
		}

		continue
	softerror:
		log.Debug().
			Err(err).
			Any("addr", conn.RemoteAddr()).
			Msg("Websocket error")
	}

error:
	if err != nil {
		return
	}

	if _, ok := err.(*websocket.CloseError); ok {
		return
	}
	log.Debug().
		Err(err).
		Any("addr", conn.RemoteAddr()).
		Msg("Websocket fatal error (closing connection)")
	conn.Close()
}

func (viewer *ViewerServer) Start() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60 * time.Second))

	router.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warn().Err(err).Msg("Upgrade failed")
			return
		}

		go viewer.handleClient(conn)
	})

	log.Info().Str("addr", viewer.Addr).Msg("Viewer server starting")
	http.ListenAndServe(viewer.Addr, router)
}
