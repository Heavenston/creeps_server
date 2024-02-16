package viewer

import (
	"embed"
	"io/fs"
	"net/http"
	"time"

	"creeps.heav.fr/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ViewerServer struct {
    Server *server.Server
    Addr string
}

func (viewer *ViewerServer) handleClient(conn *websocket.Conn) {
    defer conn.Close()
    defer log.Debug().Any("addr", conn.RemoteAddr()).Msg("Websocket connection closed")

    log.Debug().Any("addr", conn.RemoteAddr()).Msg("New websocket connection")
    for {
        t, _, err := conn.NextReader()
        if _, ok := err.(*websocket.CloseError); ok {
            return
        }
        if err != nil {
            log.Debug().Err(err).Any("addr", conn.RemoteAddr()).Msg("Websocket read error")
            return
        }

        if t != websocket.TextMessage {
            log.Debug().Any("addr", conn.RemoteAddr()).Msg("Non text websocket message (colsed connection)")
            return
        }
    }
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

