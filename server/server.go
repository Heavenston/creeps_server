package server

import (
	"creeps.heav.fr/server/model"
	"creeps.heav.fr/server/terrain"
)

type Server struct {
	tilemap *terrain.Tilemap
	ticker  *Ticker

	setup *model.SetupResponse
	costs *model.CostsResponse
}

func NewServer(tilemap *terrain.Tilemap, setup *model.SetupResponse, costs *model.CostsResponse) *Server {
	srv := new(Server)
	srv.tilemap = tilemap
	srv.ticker = NewTicker(setup.TicksPerSeconds)
	srv.setup = setup
	srv.costs = costs
	return srv
}

func (srv *Server) Ticker() *Ticker {
	return srv.ticker
}

func (srv *Server) Tilemap() *terrain.Tilemap {
	return srv.tilemap
}

func (srv *Server) Start() {
	srv.ticker.Start()
}
