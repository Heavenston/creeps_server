package api

import (
	"encoding/json"
	"errors"
	"fmt"
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
    w.WriteHeader(200)
    w.Write(make([]byte, 0))

    login := chi.URLParam(r, "login")
    unitIdStr := chi.URLParam(r, "unitId")
    unitId := uid.Uid(unitIdStr)
    opcode := chi.URLParam(r, "opcode")

    sendError := func(code string, mess string) {
        bytes, err := json.Marshal(model.CommandResponse {
            OpCode: opcode,
            Login: login,
            ErrorCode: &code,
            Error: &mess,
        })
        errors.Unwrap(err)
        w.Write(bytes)
        log.Trace().
            Str("login", login).Str("unitId", unitIdStr).Str("opcode", opcode).
            Str("code", code).Str("mess", mess).
            Msg("Command failed")
    }

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
    newAction.ReportId = uid.GenUid() 
    newAction.OpCode = server.ActionOpCode(opcode)

    reportIdStr := string(newAction.ReportId)
    err := unit.StartAction(newAction)

    if err != nil {
        if _, ok := err.(server.UnitBusyError); ok {
            sendError(
                "unavailable",
                "Your unit is already doing something.",
            )
            return
        }
        if err, ok := err.(server.UnsuportedActionError); ok {
            sendError(
                "unrecognized",
                fmt.Sprintf("Ocode %s is not supported, supported actions: %v", err.Tried, err.Supported),
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
        OpCode: opcode,
        ReportId: &reportIdStr,
        Login: login,
        UnitId: &unitIdStr,
        Misses: 0, // < TODO: Count misses (put in player struct)
    }

    bytes, err := json.Marshal(response)
    errors.Unwrap(err)
    w.Write(bytes)
}
