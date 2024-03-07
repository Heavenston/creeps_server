package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"os"
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
