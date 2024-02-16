package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"creeps.heav.fr/uid"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type reportHandle struct {
    api *ApiServer
}

func (h *reportHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    reportIdStr := chi.URLParam(r, "reportId")
    reportId := uid.Uid(reportIdStr)

	log.Trace().
        Str("reportId", reportIdStr).
        Str("addr", r.RemoteAddr).
		Msg("Get report Id")

    report := h.api.Server.GetReport(reportId)
    if report == nil {
        w.WriteHeader(404)
        w.Write([]byte(fmt.Sprintf(`{
            "opcode": "noreport",
            "error": "No such reprot",
            "reportId": "%s",
        }`, reportId)))
        return
    }

    res, err := json.Marshal(report)
    errors.Unwrap(err)
    
    w.WriteHeader(200)
    w.Write(res)
}
