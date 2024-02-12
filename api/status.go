package api

import "net/http"

type statusHandle struct {
    api *ApiServer
}

func (h *statusHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))
}
