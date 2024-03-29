package model

import "github.com/heavenston/creeps_server/creeps_lib/geom"

type MessageSendParameter struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

type FireParameter struct {
	Destination geom.Point `json:"destination"`
}
