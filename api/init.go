package api

import (
	"encoding/json"
	"net/http"

	"creeps.heav.fr/api/model"
	"creeps.heav.fr/gameplay"
	"creeps.heav.fr/server"
	"github.com/go-chi/chi/v5"
)

type initHandle struct {
	api *ApiServer
}

func (h *initHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(make([]byte, 0))

	username := chi.URLParam(r, "username")

	player := server.NewPlayer(username, r.RemoteAddr)
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
}
