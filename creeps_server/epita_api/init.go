package epita_api

import (
	"encoding/json"
	"net/http"
	"strings"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_server/gameplay"
	"github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/server/entities"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type initHandle struct {
	api *ApiServer
}

func (h *initHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.RemoteAddr, ":")[0]

	username := chi.URLParam(r, "username")

	spawnPoint := h.api.Server.FindSpawnPoint(Point{}, 2, func(p Point) bool {
		found := false
		h.api.Server.ForEachEntity(func(entity server.IEntity) (shouldStop bool) {
			eplayer, ok := entity.(*entities.Player)
			if !ok {
				return
			}
			found = eplayer.GetSpawnPoint().Dist(p) < 30
			shouldStop = found
			return
		})
		return !found
	})
	player := entities.NewPlayer(h.api.Server, username, addr, spawnPoint)
	player.SetResources(h.api.Server.GetDefaultPlayerResources())
	townhall, household, c1, c2 := gameplay.InitPlayer(h.api.Server, player)

	response := model.InitResponse{}

	c1id := c1.GetId()
	response.Citizen1Id = &c1id
	c2id := c2.GetId()
	response.Citizen2Id = &c2id
	response.Costs = h.api.Server.GetCosts()
	response.Setup = h.api.Server.GetSetup()
	response.HouseholdCoordinates = &household
	response.TownHallCoordinates = &townhall
	response.Login = username
	response.PlayerId = new(int16)
	res := player.GetResources()
	response.Resources = &res
	response.Tick = h.api.Server.Ticker().GetTickNumber()

	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Write(data)

	log.Info().Str("username", username).
		Any("townhall_posiion", townhall).
		Any("c1id", c1id).
		Any("c2id", c2id).
		Any("id", player.GetId()).
		Any("addr", player.GetAddr()).
		Msg("New player init")
}
