package webserver

import (
	"net/http"
	"net/url"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	gamemanager "github.com/Heavenston/creeps_server/creeps_manager/game_manager"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var DIST_FOLDER = "./front/dist"

type WebServer struct {
    Db *gorm.DB
    GameManager *gamemanager.GameManager
    LoginURL *url.URL

    DiscordAuth *discordapi.DiscordAppAuth
}

func (self *WebServer) Start(addr string) error {
    router := chi.NewRouter()
    router.Use(self.fillCtxMiddleware)
    router.Use(self.authMiddleware)

    router.Get("/", self.getIndex)
    router.Get("/createGame", self.getCreateGame)
    router.Get("/login", self.getLogin)
    router.Get("/logout", self.getLogout)

    router.Route("/htmx", func(r chi.Router) {
        r.Use(self.htmxMiddleware)

        r.Post("/createGame", self.postHtmxCreateGame)
    })

    router.Handle("/*", http.FileServer(http.Dir(DIST_FOLDER)))

    log.Info().Str("address", addr).Msg("Starting web server")
    return http.ListenAndServe(addr, router)
}
