package server

import (
	"creeps.heav.fr/server/model"
	"creeps.heav.fr/server/terrain"
)

type UnitOwner string

type Server struct {
	tilemap *terrain.Tilemap
	ticker  *Ticker

	setup *model.SetupResponse
	costs *model.CostsResponse

	units []iunit
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

func (srv *Server) RegisterUnit(unit iunit) {
	for _, ounit := range srv.units {
		if ounit == unit {
			panic("unit already registered")
		}
	}

	unit.setId(len(srv.units))
	srv.units = append(srv.units, unit)
}

func (srv *Server) Start() {
	srv.ticker.Start()
}
