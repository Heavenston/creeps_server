package webserver

import (
	"net/http"
	"net/url"
	"time"

	creepsjwt "github.com/Heavenston/creeps_server/creeps_manager/creeps_jwt"
	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/rs/zerolog/log"
)

func (self *WebServer) getLogin(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		nurl := *self.LoginURL
		qq := nurl.Query()
		if state, ok := r.URL.Query()["state"]; ok {
			qq["state"] = state
		}
		nurl.RawQuery = qq.Encode()

		w.Header().Add("Location", nurl.String())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	atr, err := discordapi.MakeAccessTokenRequest(
		self.DiscordAuth,
		code,
		self.LoginURL.Query().Get("redirect_uri"),
	)
	if err != nil {
		w.Header().Add("Location", "/?error="+url.PathEscape("The provided discord code is not valid"))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	discordUser, err := discordapi.GetCurrentUser(&atr)
	if err != nil {
		w.Header().Add("Location", "/?error="+url.PathEscape("The provided discord code is not valid"))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	var user model.User
	self.Db.Where("discord_id = ?", discordUser.Id).
		Assign(model.User{
			DiscordId: discordUser.Id,
			DiscordAuth: model.UserDiscordAuth{
				AccessToken:  atr.AccessToken,
				TokenExpires: time.Now().Add(time.Duration(atr.ExpiresIn) * time.Second),
				RefreshToken: atr.RefreshToken,
				Scope:        atr.Scope,
			},
		}).
		FirstOrCreate(&user)

	strToken, err := creepsjwt.Encode(int(user.ID))
	if err != nil {
		log.Error().Err(err).Msg("token encode error")
		w.Header().Add("Location", "/?error="+url.PathEscape("Internal server error"))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "creeps_token",
		Value:    strToken,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   60 * 60 * 24,
	})
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (self *WebServer) getLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/")
	http.SetCookie(w, &http.Cookie{
		Name:     "creeps_token",
		Value:    "",
		MaxAge:   0,
		SameSite: http.SameSiteDefaultMode,
	})
	w.WriteHeader(http.StatusTemporaryRedirect)
}
