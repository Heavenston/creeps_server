package viewer

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"time"

	"creeps.heav.fr/epita_api/model"
	"creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
	. "creeps.heav.fr/geom"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ViewerServer struct {
    Server *server.Server
    Addr string
}

type message struct {
    Kind string `json:"kind"`
    Content json.RawMessage `json:"content"`
}

// first packet sent by the server to the client as soon as the connection is
// established with various informations
type initContent struct {
    ChunkSize int `json:"chunkSize"`
    Setup *model.SetupResponse `json:"setup"`
    Costs *model.CostsResponse `json:"costs"`
}

// sent by the front end to subscribe to a chunk content
type subscribeRequestContent struct {
    ChunkPos Point `json:"chunkPos"`
}

// sent by the front end to unsubscribe from a chunk content
type unsubscribeRequestContent struct {
    ChunkPos Point `json:"chunkPos"`
}

func (viewer *ViewerServer) handleClient(conn *websocket.Conn) {
    var err error = nil

    defer conn.Close()
    defer log.Debug().Any("addr", conn.RemoteAddr()).Msg("Websocket connection closed")

    log.Debug().Any("addr", conn.RemoteAddr()).Msg("New websocket connection")

    {
        var initMessage message
        initMessage.Kind = "init"
        var messContent initContent
        messContent.ChunkSize = terrain.ChunkSize
        messContent.Costs = viewer.Server.GetCosts()
        messContent.Setup = viewer.Server.GetSetup()

        bytes, err := json.Marshal(messContent)
        if err != nil { goto error }
        initMessage.Content = bytes

        bytes, err = json.Marshal(initMessage)
        if err != nil { goto error }

        err = conn.WriteMessage(websocket.TextMessage, bytes)
        if err != nil { goto error }
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
        default:
            log.Debug().
                Any("addr", conn.RemoteAddr()).
                Any("mess", mess).
                Any("kind", mess.Kind).
                Msg("Unknown message")
        }
    }

error:
    if _, ok := err.(*websocket.CloseError); ok {
        return
    }
    log.Debug().Err(err).Any("addr", conn.RemoteAddr()).Msg("Websocket error")
}

//go:embed front/dist
var frontFs embed.FS

func (viewer *ViewerServer) Start() {
    upgrader := websocket.Upgrader{
    	ReadBufferSize:  1024,
    	WriteBufferSize: 1024,
    }

    router := chi.NewRouter()
    router.Use(middleware.RealIP)
    router.Use(middleware.Recoverer)
    router.Use(middleware.Timeout(60 * time.Second))

    subFs, err := fs.Sub(frontFs, "front/dist")
    if err != nil {
        log.Fatal().Err(err)
    }
    router.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Warn().Err(err).Msg("Upgrade failed")
            return
        }

        go viewer.handleClient(conn)
    })

    router.Handle("/*", http.FileServer(http.FS(subFs)))

    log.Info().Str("addr", viewer.Addr).Msg("Viewer server starting")
    http.ListenAndServe(viewer.Addr, router)
}

