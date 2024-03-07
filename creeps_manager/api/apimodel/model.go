package apimodel

import (
	"fmt"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
)

type User struct {
	Id         int     `json:"id"`
	DiscordId  string  `json:"discord_id"`
	DiscordTag string  `json:"discord_tag"`
	AvatarUrl  *string `json:"avatar_url"`
	Username   string  `json:"username"`
}

func UserFromModel(user model.User) (result User, err error) {
	discordUser, err := discordapi.GetCurrentUser(&discordapi.DiscordBearerAuth{
		DiscordId:   &user.DiscordId,
		AccessToken: user.DiscordAuth.AccessToken,
	})
	if err != nil {
		return
	}

	result = User{
		Id:         int(user.ID),
		DiscordId:  user.DiscordId,
		DiscordTag: discordUser.Discriminator,
		AvatarUrl:  nil,
		Username:   discordUser.Username,
	}

	if discordUser.Avatar != nil {
		url := fmt.Sprintf(
			"https://cdn.discordapp.com/avatars/%s/%s.png",
			discordUser.Id, *discordUser.Avatar,
		)
		result.AvatarUrl = &url
	}

	return
}

type GameConfig struct {
	CanJoinAfterStart bool `json:"can_join_after_start"`
	Private           bool `json:"private"`
	IsLocal           bool `json:"is_local"`
}

type Game struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Config GameConfig `json:"config"`

	Creator User   `json:"creator"`
	Players []User `json:"players"`

	StartedAt *int64 `json:"started_at,omitempty"`
	EndedAt   *int64 `json:"ended_at,omitempty"`
}

func GameFromModel(game model.Game) (result Game, err error) {
	result.Id = int(game.ID)
	result.Name = game.Name
	result.Config = GameConfig(game.Config)

	result.Creator, err = UserFromModel(*game.Creator)
	if err != nil {
		result = Game{}
		return
	}

	result.Players = []User{}
	for _, player := range game.Players {
		var p User
		p, err = UserFromModel(player)
		if err != nil {
			result = Game{}
			return
		}
		result.Players = append(result.Players, p)
	}

	if game.StartedAt != nil {
		tm := game.StartedAt.Unix()
		result.StartedAt = &tm
	}
	if game.EndedAt != nil {
		tm := game.EndedAt.Unix()
		result.EndedAt = &tm
	}

	return
}
