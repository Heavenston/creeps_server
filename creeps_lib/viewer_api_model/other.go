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

// Sent as a response to C2SInit
type ActionData struct {
	ActionOpCode model.ActionOpCode `json:"actionOpCode"`
	ReportId     uid.Uid            `json:"reportId"`
	Parameter    any                `json:"parameter,omitempty"`
}

// sent by the server when a unit dies or gets out of the subscribed chunks
type PlayerDespawnContent struct {
	Id uid.Uid `json:"id"`
}
