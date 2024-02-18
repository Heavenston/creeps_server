package server

import (
	"fmt"
	"math/rand"
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/server/terrain"
	"creeps.heav.fr/spatialmap"
	"creeps.heav.fr/uid"
	"github.com/rs/zerolog/log"
)

type Server struct {
	tilemap *terrain.Tilemap
	ticker  *Ticker

	setup *model.SetupResponse
	costs *model.CostsResponse

	units       *spatialmap.SpatialMap[IUnit]
	playersLock sync.RWMutex
	players     map[uid.Uid]*Player
	reportsLock sync.RWMutex
	// reports older than (somevalue) gets removed by the creeps garbage collector
	reports     map[uid.Uid]model.IReport

	defaultPlayerResourcesLock sync.RWMutex
	defaultPlayerResources     model.Resources

	randLock  sync.Mutex
	spawnRand rand.Rand
}

func NewServer(tilemap *terrain.Tilemap, setup *model.SetupResponse, costs *model.CostsResponse) *Server {
	srv := new(Server)
	srv.tilemap = tilemap

	srv.units = spatialmap.Make[IUnit]()
	srv.players = make(map[uid.Uid]*Player)
	srv.reports = make(map[uid.Uid]model.IReport)

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

// returns the spatial map with all units, do NOT use it to add or remove
// units, only use Server.RegisterUnit and Server.RemoveUnit for that
func (srv *Server) Units() *spatialmap.SpatialMap[IUnit] {
	return srv.units
}

func (srv *Server) RegisterUnit(unit IUnit) {
	if unit.GetServer() != srv {
		panic("Cannot register unit made for another server")
	}

	srv.units.Add(unit)
}

func (srv *Server) RemoveUnit(id uid.Uid) IUnit {
	unit := srv.units.RemoveFirst(func(unit IUnit) bool {
		return unit.GetId() == id
	})
	if unit == nil {
		return nil
	}
	return *unit
}

func (srv *Server) GetUnit(id uid.Uid) IUnit {
	units := srv.units.Copy()

	next := units.Iter()
	for ok, _, unit := next(); ok; ok, _, unit = next() {
		if (*unit).GetId() == id {
			return (*unit)
		}
	}
	return nil
}

// You probably want to use gameplay.InitPlayer instead
func (srv *Server) RegisterPlayer(player *Player) {
	srv.playersLock.Lock()
	defer srv.playersLock.Unlock()

	present := srv.players[player.id]
	if present == player {
		panic("Player " + player.id + " already registred")
	} else if present != nil {
		panic("A player with id " + player.id + " is already registred")
	}
	srv.players[player.id] = player
}

func (srv *Server) RemovePlayer(id uid.Uid) *Player {
	srv.playersLock.Lock()
	defer srv.playersLock.Unlock()

	player := srv.players[id]
	if player == nil {
		return nil
	}
	delete(srv.players, id)
	return player
}

func (srv *Server) GetPlayerFromId(id uid.Uid) *Player {
	srv.playersLock.RLock()
	defer srv.playersLock.RUnlock()

	return srv.players[id]
}

func (srv *Server) GetPlayerFromUsername(username string) *Player {
	srv.playersLock.RLock()
	defer srv.playersLock.RUnlock()

	for _, player := range srv.players {
		if player.GetUsername() == username {
			return player
		}
	}
	return nil
}

func (srv *Server) AddReport(report model.IReport) {
	srv.reportsLock.Lock()
	defer srv.reportsLock.Unlock()

	if srv.reports[report.GetReport().ReportId] != nil {
		panic(fmt.Errorf("cannot add a report twice (%s)", report.GetReport().ReportId))
	}

	srv.reports[report.GetReport().ReportId] = report
}

func (srv *Server) GetReport(id uid.Uid) model.IReport {
	srv.reportsLock.Lock()
	defer srv.reportsLock.Unlock()
	return srv.reports[id]
}

func (srv *Server) GetSetup() *model.SetupResponse {
	return srv.setup
}

func (srv *Server) GetCosts() *model.CostsResponse {
	return srv.costs
}

func (srv *Server) Start() {
	log.Info().Msg("Server starting")
	srv.ticker.Start()
}

func (srv *Server) SetDefaultPlayerResources(resources model.Resources) {
	srv.defaultPlayerResourcesLock.Lock()
	defer srv.defaultPlayerResourcesLock.Unlock()
	srv.defaultPlayerResources = resources
}

func (srv *Server) GetDefaultPlayerResources() model.Resources {
	srv.defaultPlayerResourcesLock.RLock()
	defer srv.defaultPlayerResourcesLock.RUnlock()
	return srv.defaultPlayerResources
}

func (srv *Server) emptyProportion(point Point) (bool, float64, float64) {
	count := 0.
	grass_count := 0
	sum_x := 0.
	sum_y := 0.

	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			np := point.Plus(dx, dy)

			tile := srv.tilemap.GetTile(np)
			if tile.Kind == terrain.TileGrass {
				sum_x += float64(np.X)
				sum_y += float64(np.X)
				grass_count += 1
			}

			count++
		}
	}

	return int(count) == grass_count, sum_x / count, sum_y / count
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

	dist := 5

	srv.randLock.Lock()
	defer srv.randLock.Unlock()

	log.Trace().Int("dist", dist).Msg("[SPAWN_POINT] Looking for spawn point...")
	for dist < 1_000_000_000 {
		for try := 0; try < 120; try++ {
			center := Point{X: srv.spawnRand.Intn(dist*2) - dist, Y: srv.spawnRand.Intn(dist*2) - dist}

			point, found := srv.FindSpawnPointNear(center)
			if found {
				log.Trace().Any("point", point).Msg("[SPAWN_POINT] Found")
				return point
			}
		}
		log.Trace().Int("dist", dist).Msg("[SPAWN_POINT] not found, increasing dist")
		dist += dist / 2
	}

	panic("could not find spawn point")
}
