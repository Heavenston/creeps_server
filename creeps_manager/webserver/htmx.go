package webserver

import (
	"net/http"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/rs/zerolog/log"
)

func (self *WebServer) postHtmxCreateGame(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        w.Write([]byte(`<popup-spawn kind="error">Bad request</popup-spawn>`))
        return
    }

    user, ok := r.Context().Value("user").(model.User)
    if !ok {
        w.Write([]byte(`<popup-spawn kind="error">You are not logged in</popup-spawn>`))
        return
    }

    gameName := r.FormValue("game_name")
    if gameName == "" {
        w.Write([]byte(`<popup-spawn kind="error">Invalid game name</popup-spawn>`))
        return
    }

    var count int64
    rs := self.Db.Model(&model.Game{}).
        Where("creator_id = ?", user.ID).
        Where("ended_at <> NULL").
        Count(&count)
    if rs.Error != nil {
        log.Error().Err(rs.Error).Msg("DB Error")
        w.Write([]byte(`<popup-spawn kind="error">Internal Server Error</popup-spawn>`))
        return
    }

    if count >= 5 {
        w.Write([]byte(`<popup-spawn kind="error">You cannot create more than 5 games</popup-spawn>`))
        return
    }

    game := model.Game {
        Name: gameName,
        CreatorID: int(user.ID),
        Config: model.GameConfig{
            CanJoinAfterStart: true,
            Private: false,
            IsLocal: false,
        },
    }

    rs = self.Db.Create(&game)
    if rs.Error != nil {
        log.Error().Err(rs.Error).Msg("DB Error")
        w.Write([]byte(`<popup-spawn kind="error">Internal Server Error</popup-spawn>`))
        return
    }

    log.Info().Str("name", game.Name).Uint("id", game.ID).Msg("Created game")

    _, err = self.GameManager.StartGame(game)
    if err != nil {
        log.Error().Err(rs.Error).Msg("Start game error")
        w.Write([]byte(`<popup-spawn kind="error">
            Internal server error while starting the game
        </popup-spawn>`))
        return
    }

    w.Header().Add("HX-Location", "/")
}

