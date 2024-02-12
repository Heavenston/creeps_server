package api

import (
	"fmt"
	"net/http"
)

type reportHandle struct {
    api *ApiServer
}

func (h *reportHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))

    fmt.Printf("Ola %s", r.URL)
}
