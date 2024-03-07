package api

import (
	"encoding/json"
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
		w.Write([]byte(`{"erro":"internal_error", "message": "internal error"}`))
		return
	}

	body, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("serialization error")
		w.WriteHeader(500)
		w.Write([]byte(`{"erro":"internal_error", "message": "internal error"}`))
		return
	}

	w.Header().Add("Content-Type","application/json")
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
			w.Write([]byte(`{"erro":"internal_error", "message": "internal error"}`))
			return
		}
		results = append(results, result)
	}

	body, err := json.Marshal(results)
	if err != nil {
		log.Error().Err(err).Msg("serialization error")
		w.WriteHeader(500)
		w.Write([]byte(`{"erro":"internal_error", "message": "internal error"}`))
		return
	}

	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(200)
	w.Write(body)
}
