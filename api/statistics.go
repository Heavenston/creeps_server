package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type statisticsHandle struct {
    api *ApiServer
}

func (h *statisticsHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))

    log.Trace().
        Str("addr", r.RemoteAddr).
        Msg("Statistics request")
}
