package viewer

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	. "creeps.heav.fr/geom"
	"creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
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
	socketLock       sync.Mutex
	socket           *websocket.Conn
	chunksLock       sync.RWMutex
	subscribedChunks map[Point]bool
}

func (viewer *ViewerServer) handleClientSubscription(
	chunkPos Point,
	conn *connection,
) {
	chunk := viewer.Server.Tilemap().GetChunk(chunkPos)
	if chunk == nil {
		// FIXME: DO NOT GENERATE, maybe create the object but still to-be-generated
		//        as we still need to be notified when it is generated
		chunk = viewer.Server.Tilemap().GenerateChunk(chunkPos)
	}

	terrainChangeChannel := make(chan terrain.TilemapUpdateEvent, 64)
	terrainCancelHandle := chunk.UpdatedEventProvider.Subscribe(terrainChangeChannel)
	defer terrainCancelHandle.Cancel()

	serverEventsChannel := make(chan server.IServerEvent, 64)
	aabb := AABB {
		From: chunkPos.Times(terrain.ChunkSize),
		Size: Point {
			X: terrain.ChunkSize,
			Y: terrain.ChunkSize,
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

	sendUnit := func(unit server.IUnit) {
		sendMessage("unit", unitContent{
			OpCode: unit.GetOpCode(),
			UnitId: unit.GetId(),
			Owner: unit.GetOwner(),
			Position: unit.GetPosition(),
			Upgraded: unit.IsUpgraded(),
		})
	}

	sendTerrain()

	units := viewer.Server.Units().GetAllIntersects(aabb)
	for _, unit := range units {
		sendUnit(unit)
	}

	viewer.Server.ForEachPlayer(func(player *server.Player) {
		sendMessage("playerSpawn", playerSpawnContent {
			Id: player.GetId(),
			SpawnPosition: player.GetSpawnPoint(),
			Username: player.GetUsername(),
			Resources: player.GetResources(),
		})
	})

	for {
		conn.chunksLock.RLock()
		stillSubed := conn.subscribedChunks[chunkPos]
		conn.chunksLock.RUnlock()
		if !stillSubed {
			break
		}

		select {
		case change, ok := (<-terrainChangeChannel):
			if !ok {
				log.Trace().Msg("terrain channel closed")
				break
			}

			sendMessage("tileChange", tileChangeContent {
				TilePos: change.UpdatedPosition.Add(chunkPos.Times(terrain.ChunkSize)),
				Kind: byte(change.NewValue.Kind),
				Value: change.NewValue.Value,
			})
		case event, ok := (<-serverEventsChannel):
			if !ok {
				log.Trace().Msg("server events channel closed")
				break
			}

			if e, ok := event.(*server.UnitSpawnEvent); ok {
				sendUnit(e.Unit)
			}
			if e, ok := event.(*server.UnitDespawnEvent); ok {
				sendMessage("unitDespawned", unitDespawnContent {
					UnitId: e.Unit.GetId(),
				})
			}
			if e, ok := event.(*server.UnitMovedEvent); ok {
				sendMessage("unitMovement", unitMovementContent {
					UnitId: e.Unit.GetId(),
					New: e.To,
				})
			}
			if e, ok := event.(*server.UnitUpgradedEvent); ok {
				sendMessage("unitUpgraded", unitUpgradedContent {
					UnitId: e.Unit.GetId(),
				})
			}
			if e, ok := event.(*server.PlayerSpawnEvent); ok {
				sendMessage("playerSpawn", playerSpawnContent {
					Id: e.Player.GetId(),
					SpawnPosition: e.Player.GetSpawnPoint(),
					Username: e.Player.GetUsername(),
					Resources: e.Player.GetResources(),
				})
			}
			if e, ok := event.(*server.PlayerDespawnEvent); ok {
				sendMessage("playerDespawn", playerDespawnContent {
					Id: e.Player.GetId(),
				})
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
			Any("message", mess).
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
		CheckOrigin: func(r *http.Request) bool { return true; },
	}

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
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
