package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	router.Route("/users/{userId}", func(r chi.Router) {
		r.Get("/", (&getUserHandle{cfg: cfg}).ServeHTTP)
	})
	router.Route("/games", func(r chi.Router) {
		r.Route("/{gameId}", func(r chi.Router) {
			r.Get("/", (&getGameHandle{cfg: cfg}).ServeHTTP)
		})
		r.Get("/", (&getGamesHandle{cfg: cfg}).ServeHTTP)
	})

	return router
}

func Start(cfg ApiCfg) error {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}))
	router.Use(middleware.RealIP)

	router.Mount("/api", apiRouter(&cfg))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(404)
		fmt.Fprintf(w, `{"error": "not_found", "message": "This endpoint does not exist"}`)
	})

	log.Info().Str("addr", cfg.TargetAddr).Msg("Api server starting")

	return http.ListenAndServe(cfg.TargetAddr, router)
}
