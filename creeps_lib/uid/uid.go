package uid

import (
	"math/rand"
	"strings"
)

type Uid string

// invalid Uid so impossible to have otherwise
const ServerUid Uid = "_server"

var alphabet string = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"

func GenUid() Uid {
	result := strings.Builder{}

	for i := 0; i < 10; i++ {
		result.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	return Uid(result.String())
}
