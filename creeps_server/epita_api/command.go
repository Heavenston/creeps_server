package epita_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/server/entities"
	"github.com/rs/zerolog/log"
)

type commandHandle struct {
	api *ApiServer
}

func (h *commandHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.RemoteAddr, ":")[0]

	login := chi.URLParam(r, "login")
	unitIdStr := chi.URLParam(r, "unitId")
	unitId := uid.Uid(unitIdStr)
	strOpcode := chi.URLParam(r, "opcode")
	opcode := model.ActionOpCode(strOpcode)

	sendError := func(code string, mess string) {
		bytes, err := json.Marshal(model.CommandResponse{
			OpCode:    opcode,
			Login:     login,
			UnitId:    &unitId,
			ReportId:  nil,
			ErrorCode: &code,
			Error:     &mess,
		})
		errors.Unwrap(err)
		w.Write(bytes)
		log.Trace().
			Str("login", login).Str("unitId", unitIdStr).Str("opcode", strOpcode).
			Str("code", code).Str("mess", mess).
			Msg("Command failed")
	}

	log.Debug().
		Str("login", login).Str("unitId", unitIdStr).Str("opcode", strOpcode).
		Msg("Command post")

	if !opcode.IsValid() {
		sendError(
			"unrecognized",
			fmt.Sprintf("Opcode '%s' doesn't exist", opcode),
		)
		return
	}

	player, _ := h.api.Server.FindEntity(func(e server.IEntity) bool {
		if p, ok := e.(*entities.Player); ok {
			return p.GetUsername() == login
		}
		return false
	}).(*entities.Player)

	if player == nil || player.GetAddr() != addr {
		log.Trace().
			Str("login", login).
			Bool("found", player != nil).
			Str("addr", addr).
			Msg("Access denied")

		sendError(
			"noplayer",
			"The login you provided does not exist or is not someone you have access to",
		)
		return
	}

	unit, _ := h.api.Server.GetEntity(unitId).(server.IUnit)

	if unit == nil || unit.GetOwner() != player.GetId() {
		sendError(
			"nounit",
			"The unitId you provided did not match any of your units.",
		)
		return
	}

	if !unit.IsRegistered() {
		sendError("dead", "Your unit died.")
		return
	}

	newAction := new(server.Action)
	newAction.ReportId = uid.GenUid()
	newAction.OpCode = opcode

	paramType := opcode.ParameterType()
	if paramType != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sendError("invalidparameter", "Cannot read the body")
			log.Warn().Err(err).Msg("Error while reading request body")
			return
		}

		paramValue := reflect.New(paramType)
		err = json.Unmarshal(body, paramValue.Interface())
		if err != nil {
			sendError("invalidparameter", "Cannot deserialize the body")
			return
		}
		newAction.Parameter = reflect.Indirect(paramValue).Interface()
	}

	err := unit.StartAction(newAction, nil)

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
	}

	response := model.CommandResponse{
		OpCode:   opcode,
		ReportId: &newAction.ReportId,
		Login:    login,
		UnitId:   &unitId,
		Misses:   0, // < TODO: Count misses (put in player struct)
	}

	bytes, err := json.Marshal(response)
	errors.Unwrap(err)
	w.Write(bytes)
}
