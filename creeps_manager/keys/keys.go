package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
)

var PrivKey *rsa.PrivateKey
var PubKey *rsa.PublicKey

func init() {
    var err error
    PrivKey, err = rsa.GenerateKey(rand.Reader, 256)
    if err != nil {
        fmt.Printf("Could not generate rsa private key: %e", err)
        os.Exit(1)
    }
    PubKey = &PrivKey.PublicKey
}
