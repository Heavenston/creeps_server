package epita_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "creeps.heav.fr/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
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
    router := chi.NewRouter()
    router.Use(middleware.RealIP)
    router.Use(middleware.Recoverer)
    router.Use(middleware.Timeout(60 * time.Second))

    router.Handle("/status", &statusHandle {
        api: api,
    })
    
    router.Handle("/statistics", &statisticsHandle {
        api: api,
    })
    
    router.Handle("/init/{username}", &initHandle {
        api: api,
    })
    
    router.Handle("/command/{login}/{unitId}/{opcode}", &commandHandle {
        api: api,
    })
    
    router.Handle("/report/{reportId}", &reportHandle {
        api: api,
    })
    
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(404)
        fmt.Printf("not found: %s %s\n", r.Method, r.URL)
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

    log.Info().Str("addr", api.Addr).Msg("Api server starting")
    
    http.ListenAndServe(api.Addr, router)
}
