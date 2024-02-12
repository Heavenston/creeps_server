package server

import (
	"math/rand"
	"strings"
)

type Uid string

const ServerUid Uid = "_server"

var alphabet string = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"

func GenUid() Uid {
	result := strings.Builder{}

	for i := 0; i < 5; i++ {
		result.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	return Uid(result.String())
}
