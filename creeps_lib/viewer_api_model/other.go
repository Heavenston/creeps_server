package viewerapimodel

import (
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
)

//tygo:emit
var _ = `
export type Message =
| C2SMessage
| S2CMessage
;
`

type ActionData struct {
	ActionOpCode model.ActionOpCode `json:"actionOpCode"`
	ReportId     uid.Uid            `json:"reportId"`
	Parameter    any                `json:"parameter,omitempty"`
}
