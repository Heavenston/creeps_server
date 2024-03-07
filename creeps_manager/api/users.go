package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Heavenston/creeps_server/creeps_manager/api/apimodel"
	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type usersHandle struct {
	cfg *ApiCfg
}

func (h *usersHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    userId := chi.URLParam(r, "userId")

    var user model.User

    if userId == "@me" {
        var err error
        user, err = auth(h.cfg.Db, w, r)
        if err != nil {
            return
        }
    } else {
        rs := h.cfg.Db.Where("ID = ?", userId).First(&user)
        if rs.RowsAffected == 0 {
            w.Header().Add("content-type", "application/json")
            w.WriteHeader(404)
            w.Write([]byte(`{"error": "not_found", "message": "Could not find the user with the given userId"}`))
        }
    }

    discordUser, err := discordapi.GetCurrentUser(&discordapi.DiscordBearerAuth{
        DiscordId: &user.DiscordId,
        AccessToken: user.DiscordAuth.AccessToken,
    })
    if err != nil {
        return
    }

    result := apimodel.User{
    	Id: int(user.ID),
    	DiscordId: user.DiscordId,
    	DiscordTag: discordUser.Discriminator,
        AvatarUrl: nil,
    	Username: discordUser.Username,
    }

    if discordUser.Avatar != nil {
        url := fmt.Sprintf(
            "https://cdn.discordapp.com/avatars/%s/%s.png",
            discordUser.Id, *discordUser.Avatar,
        )
        result.AvatarUrl = &url
    }

    data, err := json.Marshal(result)
    if err != nil {
        log.Warn().Err(err).Msg("serialization error")
        w.WriteHeader(500)
        w.Write([]byte(`{"erro":"internal_error", "message": "internal error"}`))
        return
    }

    w.Header().Add("content-type", "application/json")
    w.WriteHeader(200)
    w.Write(data)
}
