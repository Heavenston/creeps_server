package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ApiCfg struct {
	Db         *gorm.DB
	TargetAddr string

	DiscordAuth *discordapi.DiscordAppAuth
}

func apiRouter(cfg *ApiCfg) http.Handler {
	router := chi.NewRouter()

	router.Get("/login", (&loginHandle{cfg: cfg}).ServeHTTP)

	return router
}

func Start(cfg ApiCfg) error {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Mount("/api", apiRouter(&cfg))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(404)
		fmt.Fprintf(w, `{"error": "not_found", "message": "This endpoint does not exist"}`)
	})

	log.Info().Str("addr", cfg.TargetAddr).Msg("Api server starting")

	return http.ListenAndServe(cfg.TargetAddr, router)
}
