package viewerapimodel

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type Message struct {
	Kind    string          `json:"kind"`
	Content json.RawMessage `json:"content"`
}

type IMsgContent interface {
	MsgKind() string
} 

func CreateMessage(content IMsgContent) Message {
	contentbytes, err := json.Marshal(content)
	if err != nil {
		log.Fatal().Err(err).Msg("full chunk ser error")
	}

	return Message {
		Kind: content.MsgKind(),
		Content: contentbytes,
	}
}
