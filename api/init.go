package api

import "net/http"

type initHandle struct {
    api *ApiServer
}

func (h *initHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))
}
