package model

import "creeps.heav.fr/geom"

type MessageSendParameter struct {
    Recipient string `json:"recipient"`
    Message string `json:"message"`
}

type FireParameter struct {
    Destination geom.Point `json:"destination"`
}
