package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"creeps.heav.fr/api/model"
	"creeps.heav.fr/server"
	"creeps.heav.fr/uid"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type commandHandle struct {
    api *ApiServer
}

func (h *commandHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    sendError := func(code string, mess string) {
        bytes, err := json.Marshal(model.CommandResponse {
            ErrorCode: &code,
            Error: &mess,
        })
        errors.Unwrap(err)
        w.Write(bytes)
    }

    w.WriteHeader(200)
    w.Write(make([]byte, 0))

    login := chi.URLParam(r, "login")
    unitIdStr := chi.URLParam(r, "unitId")
    unitId := uid.Uid(unitIdStr)
    opcode := chi.URLParam(r, "opcode")

    log.Debug().
        Str("login", login).Str("unitId", unitIdStr).Str("opcode", opcode).
        Msg("Command post")

    player := h.api.Server.GetPlayerFromUsername(login)

    if player == nil || player.GetAddr() != r.RemoteAddr {
        sendError(
            "noplayer",
            "The login you provided does not exist or is not someone you have access to",
        )
        return
    }

    unit := h.api.Server.GetUnit(unitId)

    if unit == nil || unit.GetOwner() != player.GetId() {
        errorcode := "nounit"
        error := "The unitId you provided did not match any of your units."

        bytes, err := json.Marshal(model.CommandResponse {
            ErrorCode: &errorcode,
            Error: &error,
        })
        errors.Unwrap(err)
        w.Write(bytes)
        return
    }

    if !unit.GetAlive() {
        sendError("dead", "Your unit died.")
        return
    }

    newAction := new(server.Action)
    reportIdStr := string(newAction.ReportId)
    newAction.ReportId = uid.GenUid() 

    err := unit.StartAction(newAction)

    if err != nil {
        if _, ok := err.(server.UnitBusyError); ok {
            sendError(
                "unavailable",
                "Your unit is already doing something.",
            )
            return
        }
        if _, ok := err.(server.UnsuportedActionError); ok {
            sendError(
                "unrecognized",
                "The opcode you requested was not recognized by the unit.",
            )
            return
        }
        if _, ok := err.(server.NotEnoughResourcesError); ok {
            sendError(
                "noresources",
                "You do not own enough resources at the moment. Try again later when you do.",
            )
            return
        }
    }

    response := model.CommandResponse {
        OpCode: &opcode,
        ReportId: &reportIdStr,
        Login: &login,
        UnitId: &unitIdStr,
        Misses: 0, // < TODO: Count misses (put in player struct)
    }

    bytes, err := json.Marshal(response)
    errors.Unwrap(err)
    w.Write(bytes)
}
