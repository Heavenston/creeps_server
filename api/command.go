package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type commandHandle struct {
    api *ApiServer
}

func (h *commandHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))

    login := chi.URLParam(r, "login")
    unitId := chi.URLParam(r, "unitId")
    opcode := chi.URLParam(r, "opcode")

    log.Debug().
        Str("login", login).Str("unitId", unitId).Str("opcode", opcode).
        Msg("Command post")
}
