package server

import (
	"sync/atomic"

	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

// every value is read only except
type Action struct {
	OpCode        model.ActionOpCode
	StartedAtTick int
	ReportId      uid.Uid
	Finised       atomic.Bool
	// Contains the type returned by OpCode.ParameterType()
	Parameter any
}
