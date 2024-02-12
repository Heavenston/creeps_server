package server

import (
	"math"
	"math/rand"

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

	units   spatialmap.SpatialMap[IUnit]
	players map[Uid]*Player

	spawnRand rand.Rand
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

	srv.spawnRand = *rand.New(rand.NewSource(256))
	
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

func (srv *Server) emptyProportion(point Point) (bool, float64, float64) {
	count := 0.
	sum_x := 0.
	sum_y := 0.

	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			np := point.Plus(dx, dy)

			if srv.tilemap.GetTile(np).Kind == terrain.TileGrass {
				sum_x += float64(np.X)
				sum_y += float64(np.X)
			}

			count += 1
		}
	}

	return sum_x == count, sum_x / count, sum_y / count
}

// Returns a point safe for spawning near the given point and true
// or 0,0 and false it none could be found
func (srv *Server) FindSpawnPointNear(from Point) (Point, bool) {
	visited := make(map[Point]bool)

	for {
		visited[from] = true

		all, a_x, a_y := srv.emptyProportion(from)
		if all {
			return from, true
		}

		np := Point{X: int(a_x), Y: int(a_y)}
		if visited[np] {
			if a_x > a_y {
				np.X++
			} else {
				np.Y++
			}
		}
		if visited[np] {
			return Point{}, false
		}

		from = np
	}
}

// Returns a safe spawn point
func (srv *Server) FindSpawnPoint() Point {
	// TODO: Make an algorithm to maintain some player density

	dist := 20

	for dist < math.MaxInt / 2 {
		center := Point{ X: srv.spawnRand.Intn(dist*2)-dist, Y: srv.spawnRand.Intn(dist*2)-dist }

		point, found := srv.FindSpawnPointNear(center)
		if found {
			return point
		}
		
		dist *= 2
	}

	panic("could not find spawn point")
}
