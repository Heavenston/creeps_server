package server

import (
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/server/model"
	"creeps.heav.fr/server/terrain"
	"creeps.heav.fr/spatialmap"
)

type IUnit interface {
    GetServer() *Server
    GetId() Uid
	GetAlive() bool
	SetAlive(new bool)
    // the id of the owner, note: can be the server by way of ServerUid
    GetOwner() Uid
    GetPosition() Point
    SetPosition(newPos Point)
    GetLastAction() *Action
    SetLastAction(action *Action)
	// Ran each tick after being registered by the server
	// only if GetAlive returns true
    Tick()
}

type Server struct {
	tilemap *terrain.Tilemap
	ticker  *Ticker

	setup *model.SetupResponse
	costs *model.CostsResponse

	units spatialmap.SpatialMap[IUnit]
	players map[Uid]*Player
}

func NewServer(tilemap *terrain.Tilemap, setup *model.SetupResponse, costs *model.CostsResponse) *Server {
	srv := new(Server)
	srv.tilemap = tilemap

	srv.players = make(map[Uid]*Player)
	srv.ticker = NewTicker(setup.TicksPerSeconds)
	srv.ticker.AddTickFunc(func() {
		srv.tick()	
	})

	srv.setup = setup
	srv.costs = costs
	return srv
}

func (srv *Server) tick() {
	next := srv.units.Iter()
	for ok, _, unit := next(); ok; ok, _, unit = next() {
		if !(*unit).GetAlive() {
			continue
		}
		(*unit).Tick()
	}
}

func (srv *Server) Ticker() *Ticker {
	return srv.ticker
}

func (srv *Server) Tilemap() *terrain.Tilemap {
	return srv.tilemap
}

func (srv *Server) RegisterUnit(unit IUnit) {
	if unit.GetServer() != srv {
		panic("Cannot register unit made for another server")
	}
	srv.units.Add(unit)
}

func (srv *Server) RemoveUnit(id Uid) IUnit {
	unit := srv.units.RemoveFirst(func(unit IUnit) bool {
		return unit.GetId() == id
	})
	if unit == nil {
		return nil
	}
	return *unit
}

func (srv *Server) RegisterPlayer(player *Player) {
	present := srv.players[player.id]
	if present == player {
		panic("Player " + player.id + " already registred")
	} else if present != nil {
		panic("A player with id " + player.id + " is already registred")
	}
	srv.players[player.id] = player
}

func (srv *Server) RemovePlayer(id Uid) *Player {
	player := srv.players[id]
	if player == nil {
		return nil
	}
	delete(srv.players, id)
	return player
}

func (srv *Server) GetUnit(id Uid) IUnit {
	next := srv.units.Iter()
	for ok, _, unit := next(); ok; ok, _, unit = next() {
		if (*unit).GetId() == id {
			return (*unit)
		}
	}
	return nil
}

func (srv *Server) GetSetup() *model.SetupResponse {
	return srv.setup
}

func (srv *Server) GetCosts() *model.CostsResponse {
	return srv.costs
}

func (srv *Server) Start() {
	srv.ticker.Start()
}
