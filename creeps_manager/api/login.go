package api

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	"github.com/Heavenston/creeps_server/creeps_manager/keys"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type loginHandle struct {
	cfg *ApiCfg
}

func (h *loginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allowedRedirects []string = []string{
		"http://localhost:1234?token={{token}}",
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error": "invalid_request", "message": "missing code or state query param"}`)
		return
	}

	if !slices.Contains(allowedRedirects, state) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error": "invalid_redirect", "message": "this redirect is forbidden"}`)
		return
	}

	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}
	trp, err := discordapi.MakeAccessTokenRequest(
		h.cfg.DiscordAuth,
		code,
		scheme+"://"+r.Host+r.URL.Path,
	)
	if err != nil {
		log.Warn().Err(err).Msg("access token req")
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error": "invalid_request", "message": "could not make request to discord"}`)
		return
	}

	discordUser, err := discordapi.GetCurrentUser(&trp)
	if err != nil {
		log.Warn().Err(err).Msg("get user error")
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error": "invalid_request", "message": "could not make request to discord"}`)
		return
	}

	var user model.User = model.User{
		DiscordId: discordUser.Id,
	}
	h.cfg.Db.Where("discord_id", discordUser.Id).
		Assign(model.User{
			DiscordAuth: model.UserDiscordAuth{
				AccessToken:  trp.AccessToken,
				TokenExpires: time.Now().Add(time.Duration(trp.ExpiresIn) * time.Second),
				RefreshToken: trp.AccessToken,
				Scope:        trp.Scope,
			},
		}).
		FirstOrCreate(&user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"duid": user.DiscordId,
		"uid":  user.ID,
	})
	strToken, err := token.SignedString(keys.JWTSecret)
	if err != nil {
		log.Error().Err(err).Msg("JWT Sign error")
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "internal_error", "message": "An internal error has occured"}`)
		return
	}

	w.Header().Add("Location", strings.ReplaceAll(state, "{{token}}", strToken))
	w.WriteHeader(307)
}
