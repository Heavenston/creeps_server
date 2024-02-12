package api

import "net/http"

type statisticsHandle struct {
    api *ApiServer
}

func (h *statisticsHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))
}

