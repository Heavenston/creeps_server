package epita_api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	. "github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/rs/zerolog/log"
)

type ApiServer struct {
	Server *Server

	Listener net.Listener
}

type ApiErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Error     string `json:"error"`
}

func (api *ApiServer) Start(addr string, portFile *string) error {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Handle("/status", &statusHandle{
		api: api,
	})

	router.Handle("/statistics", &statisticsHandle{
		api: api,
	})

	router.Handle("/init/{username}", &initHandle{
		api: api,
	})

	router.Handle("/command/{login}/{unitId}/{opcode}", &commandHandle{
		api: api,
	})

	router.Handle("/report/{reportId}", &reportHandle{
		api: api,
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Printf("not found: %s %s\n", r.Method, r.URL)
		marshalled, err := json.Marshal(ApiErrorResponse{
			ErrorCode: "notfound",
			Error:     "Api endpoint does not exist",
		})
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Internal Server Error: %s", err)
			return
		}
		fmt.Fprintf(w, "%s", marshalled)
	})

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	api.Listener = listener

	log.Info().Str("addr", listener.Addr().String()).Msg("Api server listening")

	if portFile != nil {
		fs, err := os.Create(*portFile)
		if err != nil {
			return err
		}
		_, port, _ := net.SplitHostPort(listener.Addr().String())
		fmt.Fprintf(fs, "%s", port)
		err = fs.Close()
		if err != nil {
			return err
		}
		log.Info().Str("file", *portFile).Str("port", port).Msg("Written api server port")
	}

	return http.Serve(listener, router)
}
