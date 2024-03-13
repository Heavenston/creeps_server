package templates

import (
	"context"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
)

func isLoggedIn(ctx context.Context) bool {
	_, ok := ctx.Value("user").(model.User)
    return ok
}

func isInGame(ctx context.Context, game model.Game) bool {
	user, ok := ctx.Value("user").(model.User)
	if !ok {
		return false
	}

	for _, player := range game.Players {
		if player.ID == user.ID {
			return true
		}
	}

	return false
}

func isCreator(ctx context.Context, game model.Game) bool {
	user, ok := ctx.Value("user").(model.User)
	if !ok {
		return false
	}

	return game.CreatorID == int(user.ID)
}
