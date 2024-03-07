package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Heavenston/creeps_server/creeps_manager/api/apimodel"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

type getGameHandle struct {
	cfg *ApiCfg
}

func (h *getGameHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	var game model.Game
	rs := h.cfg.Db.Where("id = ?", gameId).Preload(clause.Associations).Take(&game)
	if rs.RowsAffected == 0 {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "not_found", message: "not game has the given id"}`))
		return
	}

	result, err := apimodel.GameFromModel(game)
	if err != nil {
		log.Error().Err(err).Msg("game convert error")
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
		return
	}

	body, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("serialization error")
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

type getGamesHandle struct {
	cfg *ApiCfg
}

func (h *getGamesHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var games []model.Game
	h.cfg.Db.Preload(clause.Associations).Find(&games)

	var results []apimodel.Game = make([]apimodel.Game, 0, len(games))
	for _, game := range games {
		result, err := apimodel.GameFromModel(game)
		if err != nil {
			log.Error().Err(err).Msg("game convert error")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
			return
		}
		results = append(results, result)
	}

	body, err := json.Marshal(results)
	if err != nil {
		log.Error().Err(err).Msg("serialization error")
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

type postGameHandle struct {
	cfg *ApiCfg
}

func (h *postGameHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := auth(h.cfg.Db, w, r)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(415)
		w.Write([]byte(`{"error":"invalid_content_type", "message": "only application/json is supported"}`))
		return
	}

	var count int64
	rs := h.cfg.Db.Model(&model.Game{}).Where("creator_id = ?", user.ID).Count(&count)
	if rs.Error != nil {
		log.Error().Err(rs.Error).Msg("fetch error")
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
		return
	}

	if count > 0 {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"too_much", "message": "you already made more games than allowed"}`))
		return
	}

	var request apimodel.CreateGameRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"invalid_request", "message": "request's body is in the wrong format"}`))
		return
	}

	game := model.Game{
		Name: request.Name,

		CreatorID: int(user.ID),
	}
	if request.Config != nil {
		game.Config = model.GameConfig(*request.Config)
	} else {
		game.Config = model.GameConfig{
			CanJoinAfterStart: true,
			Private:           false,
			IsLocal:           false,
		}
	}

	rs = h.cfg.Db.Create(&game)
	if rs.Error != nil {
		log.Error().Err(rs.Error).Msg("create error")
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
		return
	}

	game.Creator = &user

	resp, err := apimodel.GameFromModel(game)
	if err != nil {
		log.Error().Err(err).Msg("convert error")
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Error().Err(err).Msg("marshal error")
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal_error", "message": "internal error"}`))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(respBody)
}
