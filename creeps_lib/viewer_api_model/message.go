package viewerapimodel

import (
	"encoding/json"
)

type Message struct {
	Kind    string          `json:"kind"`
	Content json.RawMessage `json:"content"`
}
