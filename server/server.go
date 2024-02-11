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

	units []IUnit
}

func NewServer(tilemap *terrain.Tilemap, setup *model.SetupResponse, costs *model.CostsResponse) *Server {
	srv := new(Server)
	srv.tilemap = tilemap

	srv.ticker = NewTicker(setup.TicksPerSeconds)
	srv.ticker.AddTickFunc(func() {
		srv.tick()	
	})

	srv.setup = setup
	srv.costs = costs
	return srv
}

func (srv *Server) tick() {
	for _, unit := range srv.units {
		unit.Tick()
	}
}

func (srv *Server) Ticker() *Ticker {
	return srv.ticker
}

func (srv *Server) Tilemap() *terrain.Tilemap {
	return srv.tilemap
}

func (srv *Server) RegisterUnit(unit IUnit) {
	for _, ounit := range srv.units {
		if ounit == unit {
			panic("unit already registered")
		}
	}

	srv.units = append(srv.units, unit)
}

func (srv *Server) GetUnit(id Uid) IUnit {
	for _, unit := range srv.units {
		if unit.GetId() == id {
			return unit
		}
	}
	return nil
}

func (srv *Server) Start() {
	srv.ticker.Start()
}
