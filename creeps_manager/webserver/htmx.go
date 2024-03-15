package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/Heavenston/creeps_server/creeps_manager/templates"
	"github.com/rs/zerolog/log"
)

func writePopup(w http.ResponseWriter, content string) {
	w.Header().Add("HX-Retarget", "body")
	w.Header().Add("HX-Reswap", "beforeend")
	w.Header().Add("HX-Reselect", "popup-spawn")

	w.Write([]byte(fmt.Sprintf(`<popup-spawn kind="error">%s</popup-spawn>`, content)))
}

func appendPopup(w http.ResponseWriter, content string) {
	w.Write([]byte(fmt.Sprintf(
		`<popup-spawn kind="error" hx-swap-oob="beforeend, body">%s</popup-spawn>`,
		content,
	)))
}

func (self *WebServer) postHtmxCreateGame(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writePopup(w, "Bad request")
		return
	}

	user, ok := r.Context().Value("user").(model.User)
	if !ok {
		writePopup(w, "You are not logged in")
		return
	}

	gameName := r.FormValue("game_name")
	if gameName == "" {
		writePopup(w, "Invalid game name")
		return
	}

	var count int64
	rs := self.Db.Model(&model.Game{}).
		Where("creator_id = ? AND ended_at IS NULL", user.ID).
		Count(&count)
	if rs.Error != nil {
		log.Error().Err(rs.Error).Msg("DB Error")
		writePopup(w, "Internal Server Error")
		return
	}

	if count >= 5 {
		writePopup(w, "You cannot create more than 5 games")
		return
	}

	game := model.Game{
		Name:      gameName,
		CreatorID: int(user.ID),
		Config: model.GameConfig{
			CanJoinAfterStart: true,
			Private:           false,
			IsLocal:           false,
		},
	}

	rs = self.Db.Create(&game)
	if rs.Error != nil {
		log.Error().Err(rs.Error).Msg("DB Error")
		writePopup(w, "Internal Server Error")
		return
	}

	log.Info().Str("name", game.Name).Uint("id", game.ID).Msg("Created game")

	_, err = self.GameManager.StartGame(game)
	if err != nil {
		log.Error().Err(err).Msg("Start game error")
		writePopup(w, `Internal server error while starting the game`)
		return
	}

	w.Header().Add("HX-Location", "/")
}

func (self *WebServer) postHtmxJoinGame(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writePopup(w, `Bad request`)
		return
	}

	user, ok := r.Context().Value("user").(model.User)
	if !ok {
		writePopup(w, `You are not logged in`)
		return
	}

	gameIdStr := r.URL.Query().Get("gameId")
	gameId, err := strconv.ParseInt(gameIdStr, 10, 32)
	if err != nil {
		writePopup(w, "Invalid game id")
		return
	}

	var game model.Game
	rs := self.Db.Where("id = ?", gameId).
		Preload("Players").
		Preload("Creator").
		Take(&game)
	if rs.RowsAffected == 0 {
		writePopup(w, "Could not find game")
		return
	}

	if game.EndedAt != nil {
		writePopup(w, "Cannot join ended game")
		return
	}

	self.Db.Model(&game).
		Where("id = ?", gameId).
		Association("Players").
		Append(&user)

	templates.GameSidePanel(self.GameManager.GetRunningGame(game.ID), game).
		Render(r.Context(), w)
}

func (self *WebServer) postHtmxStopGame(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writePopup(w, `Bad request`)
		return
	}

	user, ok := r.Context().Value("user").(model.User)
	if !ok {
		writePopup(w, `You are not logged in`)
		return
	}

	gameIdStr := r.URL.Query().Get("gameId")
	gameId, err := strconv.ParseInt(gameIdStr, 10, 32)
	if err != nil {
		writePopup(w, "Invalid game id")
		return
	}

	var game model.Game
	rs := self.Db.Where("id = ?", gameId).
		Preload("Players").
		Preload("Creator").
		Take(&game)
	if rs.RowsAffected == 0 {
		writePopup(w, "Could not find game")
		return
	}

	if game.StartedAt == nil {
		writePopup(w, "Game has not even started yet")
		return
	}

	if game.EndedAt != nil {
		writePopup(w, "Game is already ended")
		return
	}

	if game.CreatorID != int(user.ID) {
		writePopup(w, "Only the owner can stop the game")
		return
	}

	rg := self.GameManager.GetRunningGame(game.ID)
	if rg == nil {
		log.Warn().
			Any("game_id", game.ID).
			Msg("Found non existant running game for game stated as stared in the Db")
		writePopup(w, "Unexpected server error, game not running")
		return
	}

	err = rg.Stop()
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Game stop error")
		appendPopup(w, "Unexpected server error, but game should be stopped")
	}

	self.Db.Where("id = ?", gameId).
		Preload("Players").
		Preload("Creator").
		Take(&game)
	templates.GameSidePanel(nil, game).
		Render(r.Context(), w)
}
