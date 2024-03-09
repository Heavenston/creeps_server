package webserver

import (
	"net/http"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	gamemanager "github.com/Heavenston/creeps_server/creeps_manager/game_manager"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type WebServer struct {
    Db *gorm.DB
    GameManager *gamemanager.GameManager

    DiscordAuth *discordapi.DiscordAppAuth
}

func (self *WebServer) Start(addr string) error {
    router := chi.NewRouter()

    router.Handle("/*", http.FileServer(http.Dir("./front/dist")))

    return http.ListenAndServe(addr, router)
}
