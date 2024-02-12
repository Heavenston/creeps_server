package api

import "net/http"

type commandHandle struct {
    api *ApiServer
}

func (h *commandHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))
}

