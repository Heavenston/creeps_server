package api

import (
	"encoding/json"
	"net/http"

	"github.com/Heavenston/creeps_server/creeps_manager/api/apimodel"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type getUserHandle struct {
	cfg *ApiCfg
}

func (h *getUserHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	var user model.User

	if userId == "@me" {
		var err error
		user, err = auth(h.cfg.Db, w, r)
		if err != nil {
			return
		}
	} else {
		rs := h.cfg.Db.Where("ID = ?", userId).First(&user)
		if rs.RowsAffected == 0 {
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(404)
			w.Write([]byte(`{"error": "not_found", "message": "Could not find the user with the given userId"}`))
		}
	}

	result, err := apimodel.UserFromModel(user)
	if err != nil {
		log.Error().Err(err).Msg("user convert error")
		w.WriteHeader(500)
		w.Write([]byte(`{"erro":"internal_error", "message": "internal error"}`))
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("serialization error")
		w.WriteHeader(500)
		w.Write([]byte(`{"erro":"internal_error", "message": "internal error"}`))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
