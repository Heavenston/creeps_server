package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	"github.com/Heavenston/creeps_server/creeps_manager/keys"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ApiCfg struct {
	Db         *gorm.DB
	TargetAddr string

	DiscordAuth *discordapi.DiscordAppAuth
}

func auth(db *gorm.DB, w http.ResponseWriter, req *http.Request) (user model.User, err error) {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Missing auth header"}`))
		err = fmt.Errorf("No auth header")
		return
	}

	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return keys.JWTSecret, nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("token parse error")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Invalid auth header"}`))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Invalid auth header"}`))
		return
	}

	userId := claims["uid"]
	rs := db.Where("id = ?", userId).First(&user)
	if rs.RowsAffected == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden", "message": "Could not find user"}`))
		return
	}
	return
}

func apiRouter(cfg *ApiCfg) http.Handler {
	router := chi.NewRouter()

	router.Get("/login", (&loginHandle{cfg: cfg}).ServeHTTP)
	router.Route("/users/{userId}", func(r chi.Router) {
		r.Get("/", (&usersHandle{cfg: cfg}).ServeHTTP)
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
