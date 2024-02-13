package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type reportHandle struct {
    api *ApiServer
}

func (h *reportHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))

    reportId := chi.URLParam(r, "reportId")

	log.Trace().
        Str("reportId", reportId).
        Str("addr", r.RemoteAddr).
		Msg("Get report Id")
}
