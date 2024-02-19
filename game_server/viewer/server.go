package viewer

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
	"creeps.heav.fr/uid"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ViewerServer struct {
	Server *server.Server
	Addr   string
}

type message struct {
	Kind    string          `json:"kind"`
	Content json.RawMessage `json:"content"`
}

// first packet sent by the server to the client as soon as the connection is
// established with various informations
type initContent struct {
	ChunkSize int                  `json:"chunkSize"`
	Setup     *model.SetupResponse `json:"setup"`
	Costs     *model.CostsResponse `json:"costs"`
}

type fullChunkContent struct {
	ChunkPos Point `json:"chunkPos"`
	// will be base64 encoded
	// each tile has two bytes, one for the kind and one for its value
	// see terrain/tile.go for the correspondance
	// encoded in row-major order
	// can be empty if the chunk isn't generated
	Tiles []byte `json:"tiles"`
}

type unitContent struct {
	OpCode   string `json:"opCode"`
	UnitId   uid.Uid `json:"unitId"`
	Owner    uid.Uid `json:"owner"`
	Position Point `json:"position"`
}

// sent by the front end to subscribe to a chunk content
type subscribeRequestContent struct {
	ChunkPos Point `json:"chunkPos"`
}

// sent by the front end to unsubscribe from a chunk content
type unsubscribeRequestContent struct {
	ChunkPos Point `json:"chunkPos"`
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
		})
	}

	sendTerrain()

	units := viewer.Server.Units().GetAllIntersects(aabb)
	for _, unit := range units {
		sendUnit(unit)
	}

	for {
		conn.chunksLock.RLock()
		stillSubed := conn.subscribedChunks[chunkPos]
		conn.chunksLock.RUnlock()
		if !stillSubed {
			break
		}

		select {
		case _, ok := (<-terrainChangeChannel):
			if !ok {
				log.Trace().Msg("terrain channel closed")
				break
			}

			// TODO: Send parials chunk updates
			sendTerrain()
		case event, ok := (<-serverEventsChannel):
			if !ok {
				log.Trace().Msg("server events channel closed")
				break
			}

			if e, ok := event.(*server.UnitSpawnEvent); ok {
				sendUnit(e.Unit)
			}
			if e, ok := event.(*server.UnitMovedEvent); ok {
				sendUnit(e.Unit)
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
