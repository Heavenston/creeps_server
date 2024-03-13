package webserver

import (
	"context"
	"net/http"

	creepsjwt "github.com/Heavenston/creeps_server/creeps_manager/creeps_jwt"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
)

func (self *WebServer) fillCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "login_url", self.LoginURL))
		next.ServeHTTP(w, r)
	})
}

func (self *WebServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("creeps_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		claims, err := creepsjwt.Decode(tokenCookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "creeps_token",
				Value:  "",
				MaxAge: -1,
			})
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Add("HX-Redirect", "/")
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.Header().Add("Location", "/")
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}

		var user model.User
		rs := self.Db.Where("id = ?", claims.UserId).Take(&user)
		if rs.Error != nil || rs.RowsAffected == 0 {
			http.SetCookie(w, &http.Cookie{
				Name:   "creeps_token",
				Value:  "",
				MaxAge: -1,
			})
			w.Header().Add("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		// discordUser, err := discordapi.GetCurrentUser(&discordapi.DiscordBearerAuth{
		// 	AccessToken: user.DiscordAuth.AccessToken,
		// 	DiscordId:   &user.DiscordId,
		// })
		// if !reflect.DeepEqual(user.DiscordUser, discordUser) {
		// 	user.DiscordUser = discordUser
		// 	self.Db.Model(&user).
		// 		Updates(model.User{
		// 			Model: gorm.Model{ID: user.ID},
		// 			DiscordUser: discordUser,
		// 		})
		// }

		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "creeps_token",
				Value:  "",
				MaxAge: -1,
			})
			w.Header().Add("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "discordUser", user.DiscordUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (self *WebServer) htmxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
