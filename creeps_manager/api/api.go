package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ApiCfg struct {
    Db *gorm.DB
    TargetAddr string
}

func Start(cfg ApiCfg) error {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Printf("not found: %s %s\n", r.Method, r.URL)
	})

	log.Info().Str("addr", cfg.TargetAddr).Msg("Api server starting")

	return http.ListenAndServe(cfg.TargetAddr, router)
}
