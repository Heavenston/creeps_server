package viewer_api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
	. "github.com/heavenston/creeps_server/creeps_lib/viewer_api_model"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/rs/zerolog/log"
)

type ViewerServer struct {
	Server *server.Server

	AdminPassword string
	AdminAddrs    []string

	Listener net.Listener
}

func (viewer *ViewerServer) handleClient(conn *websocket.Conn) {
	var err error = nil

	defer conn.Close()
	defer log.Debug().Any("addr", conn.RemoteAddr()).Msg("Websocket connection closed")

	log.Debug().Any("addr", conn.RemoteAddr()).Msg("New websocket connection")

	connection := connection{
		viewer:           viewer,
		socket:           conn,
		subscribedChunks: make(map[Point]bool),
		knownUnits:       make(map[uid.Uid]bool),
		knownPlayers:     make(map[uid.Uid]bool),
	}
	connection.isConnected.Store(true)
	defer connection.isConnected.Store(false)

	// recv init message
	{
		var initMessage Message
		err = conn.ReadJSON(&initMessage)
		if err != nil {
			goto error
		}

		if initMessage.Kind != "init" {
			log.Warn().Any("addr", conn.RemoteAddr()).Msg("Invalid hanshake")
			return
		}

		var initContent C2SInit
		err := json.Unmarshal(initMessage.Content, &initContent)
		if err != nil {
			goto error
		}

		connection.isAdmin.Store(false)
		if initContent.AuthPassword != "" {
			remoteAddr, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
			if len(viewer.AdminAddrs) != 0 &&
				!slices.Contains(viewer.AdminAddrs, remoteAddr) {
				log.Warn().
					Any("addr", conn.RemoteAddr()).
					Any("given_auth", initContent.AuthPassword).
					Msg("Client with the wrong addr tried to give an admin password")
				return
			}

			if initContent.AuthPassword == viewer.AdminPassword {
				connection.isAdmin.Store(true)
			} else {
				log.Warn().
					Any("addr", conn.RemoteAddr()).
					Any("given_auth", initContent.AuthPassword).
					Msg("Wrong admin password used")
				return
			}
		}
	}

	// send init message
	{
		var initMessage Message
		initMessage.Kind = "init"
		var messContent S2CInit
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

	go connection.handleClientGlobalEvents()

	for {
		var mess Message
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
			var content C2SSubscribeRequest
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
			go connection.handleClientSubscription(content.ChunkPos)
		case "unsubscribe":
			var content C2SUnsubscribeRequest
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

func (viewer *ViewerServer) Start(addr string, portFile *string) error {
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

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	viewer.Listener = listener

	log.Info().Str("addr", listener.Addr().String()).Msg("Viewer server listening")

	if portFile != nil {
		fs, err := os.Create(*portFile)
		if err != nil {
			return err
		}
		_, port, _ := net.SplitHostPort(listener.Addr().String())
		fmt.Fprintf(fs, "%s", port)
		err = fs.Close()
		if err != nil {
			return err
		}
		log.Info().Str("file", *portFile).Str("port", port).Msg("Written api server port")
	}

	return http.Serve(listener, router)
}
