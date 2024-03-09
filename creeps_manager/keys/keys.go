package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

var PrivKey *rsa.PrivateKey
var PubKey *rsa.PublicKey

var JWTSecret []byte

func init() {
	var err error
	PrivKey, err = rsa.GenerateKey(rand.Reader, 256)
	if err != nil {
		fmt.Printf("Could not generate rsa private key: %e", err)
		os.Exit(1)
	}
	PubKey = &PrivKey.PublicKey

	JWTSecret = make([]byte, 64)
	_, err = io.ReadFull(rand.Reader, JWTSecret)
	if err != nil {
		fmt.Printf("Could not generate jwt secret: %e", err)
		os.Exit(1)
	}
}

func SetJwtSecret(val []byte) {
	JWTSecret = val
	log.Warn().Msg("JWTSecret has been overriden, please only do this for debugging purposes")
}
