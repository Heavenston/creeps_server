package webserver

import (
	"net/http"
)

func (self *WebServer) postHtmxCreateGame(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`<popup-spawn kind="error">Bad request</popup-spawn>`))
    }

    w.Header().Add("HX-Location", "/")
}

