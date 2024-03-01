package epita_api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type statusHandle struct {
	api *ApiServer
}

func (h *statusHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(make([]byte, 0))

	log.Trace().
		Str("addr", r.RemoteAddr).
		Msg("Get Status")
}
