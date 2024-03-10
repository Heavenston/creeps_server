package webserver

import (
	"context"
	"net/http"
	"net/url"
	"time"

	creepsjwt "github.com/Heavenston/creeps_server/creeps_manager/creeps_jwt"
	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	gamemanager "github.com/Heavenston/creeps_server/creeps_manager/game_manager"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/Heavenston/creeps_server/creeps_manager/templates"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var DIST_FOLDER = "./front/dist"

type WebServer struct {
    Db *gorm.DB
    GameManager *gamemanager.GameManager
    LoginURL *url.URL

    DiscordAuth *discordapi.DiscordAppAuth
}

func (self *WebServer) getIndex(w http.ResponseWriter, r *http.Request) {
    var games []model.Game
    self.Db.Find(&games)

    ctx := templ.WithChildren(r.Context(), templates.Index(games))
    templates.Layout(templates.IndexHeader()).
        Render(ctx, w)
}

func (self *WebServer) getCreateGame(w http.ResponseWriter, r *http.Request) {
    ctx := templ.WithChildren(r.Context(), templates.CreateGame())
    templates.Layout(templates.CreateGameHeader()).
        Render(ctx, w)
}

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
        Assign(model.User {
            DiscordId: discordUser.Id,
            DiscordAuth: model.UserDiscordAuth{
                AccessToken: atr.AccessToken,
                TokenExpires: time.Now().Add(time.Duration(atr.ExpiresIn) * time.Second),
                RefreshToken: atr.RefreshToken,
                Scope: atr.Scope,
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
        Name: "creeps_token",
        Value: strToken,
        SameSite: http.SameSiteDefaultMode,
        MaxAge: 60 * 60 * 24,
    })
    w.Header().Add("Location", "/")
    w.WriteHeader(http.StatusTemporaryRedirect)
}

func (self *WebServer) getLogout(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Location", "/")
    http.SetCookie(w, &http.Cookie{
        Name: "creeps_token",
        Value: "",
        MaxAge: 0,
        SameSite: http.SameSiteDefaultMode,
    })
    w.WriteHeader(http.StatusTemporaryRedirect)
}

func (self *WebServer) Start(addr string) error {
    router := chi.NewRouter()
    router.Use(func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            r = r.WithContext(context.WithValue(r.Context(), "login_url", self.LoginURL))
            h.ServeHTTP(w, r)
        })
    })
    router.Use(func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenCookie, err := r.Cookie("creeps_token")
            if err != nil {
                h.ServeHTTP(w, r)
                return
            }
            claims, err := creepsjwt.Decode(tokenCookie.Value)
            if err != nil {
                http.SetCookie(w, &http.Cookie{
                    Name: "creeps_token",
                    Value: "",
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
                    Name: "creeps_token",
                    Value: "",
                    MaxAge: -1,
                })
                w.Header().Add("Location", "/")
                w.WriteHeader(http.StatusTemporaryRedirect)
                return
            }

            discordUser, err := discordapi.GetCurrentUser(&discordapi.DiscordBearerAuth{
                AccessToken: user.DiscordAuth.AccessToken,
                DiscordId: &user.DiscordId,
            })

            if err != nil {
                http.SetCookie(w, &http.Cookie{
                    Name: "creeps_token",
                    Value: "",
                    MaxAge: -1,
                })
                w.Header().Add("Location", "/")
                w.WriteHeader(http.StatusTemporaryRedirect)
                return
            }

            ctx := context.WithValue(r.Context(), "user", user)
            ctx = context.WithValue(ctx, "discordUser", discordUser)
            h.ServeHTTP(w, r.WithContext(ctx))
        })
    })

    router.Get("/", self.getIndex)
    router.Get("/createGame", self.getCreateGame)
    router.Get("/login", self.getLogin)
    router.Get("/logout", self.getLogout)

    router.Route("/htmx", func(r chi.Router) {
        r.Use(func(h http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if r.Header.Get("HX-Request") != "true" {
                    w.WriteHeader(http.StatusBadRequest)
                    return
                }

                h.ServeHTTP(w, r)
            })
        })

    })

    router.Handle("/*", http.FileServer(http.Dir(DIST_FOLDER)))

    log.Info().Str("address", addr).Msg("Starting web server")
    return http.ListenAndServe(addr, router)
}
