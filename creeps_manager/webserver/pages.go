package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/Heavenston/creeps_server/creeps_manager/templates"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func (self *WebServer) getIndex(w http.ResponseWriter, r *http.Request) {
	var games []model.Game
	self.Db.Find(&games)

	ctx := templ.WithChildren(r.Context(), templates.Index(games))
	templates.Layout(templates.IndexHeader()).
		Render(ctx, w)
}

func (self *WebServer) getCreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := templ.WithChildren(r.Context(), templates.CreateGame())
	templates.Layout(templates.CreateGameHeader()).
		Render(ctx, w)
}

func (self *WebServer) getGame(w http.ResponseWriter, r *http.Request) {
	gameIdStr := chi.URLParam(r, "gameId")
	gameId, err := strconv.ParseInt(gameIdStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var game model.Game
	rs := self.Db.Where("id = ?", gameId).Take(&game)
	if rs.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if rs.Error != nil {
		log.Error().Err(rs.Error).Msg("DB Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rg := self.GameManager.GetRunningGame(game.ID)

	url := ""
	if rg != nil {
		url = fmt.Sprintf("ws://localhost:%d/websocket", rg.ViewerPort)
	}

	ctx := templ.WithChildren(r.Context(), templates.Game(url))
	templates.Layout(templates.GameHeader()).
		Render(ctx, w)
}
