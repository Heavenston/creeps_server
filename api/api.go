package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "creeps.heav.fr/server"
)

type ApiServer struct {
    Server *Server
    Addr string
}

type ApiErrorResponse struct {
    ErrorCode string `json:"errorCode"`
    Error string `json:"error"`
}

func (api *ApiServer) Start() {
    http.Handle("/status", &statusHandle {
        api: api,
    })
    
    http.Handle("/statistics", &statisticsHandle {
        api: api,
    })
    
    http.Handle("/init", &initHandle {
        api: api,
    })
    
    http.Handle("/command", &commandHandle {
        api: api,
    })
    
    http.Handle("/report", &reportHandle {
        api: api,
    })
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(404)
        marshalled, err := json.Marshal(ApiErrorResponse {
            ErrorCode: "notfound",
            Error: "Api endpoint does not exist",
        })
        if err != nil {
            w.WriteHeader(500)
        	fmt.Fprintf(w, "Internal Server Error: %s", err)
            return
        }
    	fmt.Fprintf(w, "%s", marshalled)
    })
    
    http.ListenAndServe(api.Addr, nil)
}
