package api

import "net/http"

type reportHandle struct {
    api *ApiServer
}

func (h *reportHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write(make([]byte, 0))
}
