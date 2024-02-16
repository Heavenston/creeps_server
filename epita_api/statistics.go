package epita_api

import (
	"encoding/json"
	"errors"
	"net/http"

	"creeps.heav.fr/epita_api/model"
	"github.com/rs/zerolog/log"
)

type statisticsHandle struct {
	api *ApiServer
}

func (h *statisticsHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := model.StatisticsResponse{
		ServerId:    h.api.Server.GetSetup().ServerId,
		GameRunning: true,
		Tick:        h.api.Server.Ticker().GetTickNumber(),
		Dimension:   h.api.Server.GetSetup().WorldDimension,
		Players:     []model.Player{},
	}

	body, err := json.Marshal(resp)
	errors.Unwrap(err)

	w.WriteHeader(200)
	w.Write(body)

	log.Trace().
		Str("addr", r.RemoteAddr).
		Msg("Statistics request")
}
