package server

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"

	"creeps.heav.fr/epita_api/model"
	"creeps.heav.fr/events/spatialevents"
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

	events *spatialevents.SpatialEventProvider[IServerEvent]

	entitiesLock       sync.RWMutex
	entitiesSpatialmap *spatialmap.SpatialMap[IEntity]
	entitiesMap        map[uid.Uid]IEntity

	// reports older than (somevalue) gets removed by the creeps garbage collector
	reports     map[uid.Uid]model.IReport
	reportsLock sync.RWMutex

	defaultPlayerResourcesLock sync.RWMutex
	defaultPlayerResources     model.Resources

	randLock  sync.Mutex
	spawnRand rand.Rand
}

func NewServer(tilemap *terrain.Tilemap, setup *model.SetupResponse, costs *model.CostsResponse) *Server {
	srv := new(Server)
	srv.tilemap = tilemap

	srv.events = spatialevents.NewSpatialEventProvider[IServerEvent]()

	srv.entitiesSpatialmap = spatialmap.NewSpatialMap[IEntity]()
	srv.entitiesMap = make(map[uid.Uid]IEntity)

	srv.reports = make(map[uid.Uid]model.IReport)

	srv.ticker = NewTicker(setup.TicksPerSeconds)
	srv.ticker.AddTickFunc(func() {
		srv.tick()
	})

	srv.setup = setup
	srv.costs = costs

	srv.spawnRand = *rand.New(rand.NewSource(256))

	go (func() {
		channel := make(chan IServerEvent)
		srv.events.Subscribe(channel, AABB{})
		for {
			event, ok := (<-channel)
			if !ok {
				break
			}

			log.Trace().
				Str("type", reflect.TypeOf(event).String()).
				Any("event", event).
				Msg("Server event")
		}
	})()

	return srv
}

func (srv *Server) tick() {
	srv.entitiesLock.Lock()
	entites := make([]IEntity, 0, len(srv.entitiesMap))
	for _, entity := range srv.entitiesMap {
		entites = append(entites, entity)
	}
	srv.entitiesLock.Unlock()

	for _, entity := range srv.entitiesMap {
		entity.Tick()
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
func (srv *Server) Entities() *spatialmap.SpatialMap[IEntity] {
	return srv.entitiesSpatialmap
}

// returns the spatial map with all units, do NOT use it to add or remove
// units, only use Server.RegisterUnit and Server.RemoveUnit for that
func (srv *Server) Events() *spatialevents.SpatialEventProvider[IServerEvent] {
	return srv.events
}

// Do not call directly, use entity.Register(), otherwise events won't be 
// emitted
func (srv *Server) RegisterEntity(entity IEntity) {
	if entity.GetServer() != srv {
		panic("Cannot register entity made for another server")
	}

	srv.entitiesLock.Lock()
	defer srv.entitiesLock.Unlock()

	if srv.entitiesMap[entity.GetId()] != nil {
		panic("cannot register two entities with the same id")
	}
	srv.entitiesSpatialmap.Add(entity)
	srv.entitiesMap[entity.GetId()] = entity

	log.Trace().Str("id", string(entity.GetId())).Msg("Registered entity")

	ownerId := entity.GetOwner()
	if ownerId == uid.ServerUid {
		return
	}

	ownerEntity := srv.entitiesMap[ownerId]
	owner, isOwner := ownerEntity.(IOwnerEntity)
	if !isOwner {
		log.Warn().
			Str("entity_id", string(entity.GetId())).
			Any("owner_id", string(ownerId)).
			Msg("registered entity has an invalid owner")
		return
	}

	owner.AddEntity(entity)
}

func (srv *Server) RemoveEntity(id uid.Uid) (entity IEntity) {
	srv.entitiesLock.Lock()
	defer srv.entitiesLock.Unlock()

	entity = srv.entitiesMap[id]
	if entity == nil {
		log.Warn().
			Str("id", string(id)).
			Str("type_name", string(reflect.TypeOf(entity).String())).
			Msg("Attempted to remove an entity that isn't registered")
		return
	}
	delete(srv.entitiesMap, id)
	srv.entitiesSpatialmap.RemoveFirst(func(e IEntity) bool {
		return e.GetId() == id
	})

	ownerId := entity.GetOwner()
	if ownerId == uid.ServerUid {
		return
	}

	ownerEntity := srv.entitiesMap[ownerId]
	if ownerEntity == nil {
		return
	}
	owner, isOwner := ownerEntity.(IOwnerEntity)
	if !isOwner {
		log.Warn().
			Str("entity_type", reflect.TypeOf(entity).String()).
			Str("entity_id", string(entity.GetId())).
			Any("owner_type", reflect.TypeOf(owner).String()).
			Any("owner_id", string(ownerId)).
			Msg("removed entity has an invalid owner")
		return
	}

	removed := owner.RemoveEntity(id)
	if removed == nil {
		log.Warn().
			Str("entity_id", string(entity.GetId())).
			Any("owner_id", string(ownerId)).
			Msg("entity's owner did not have it registered")
	}
	return
}

func (srv *Server) GetEntity(id uid.Uid) IEntity {
	srv.entitiesLock.RLock()
	defer srv.entitiesLock.RUnlock()
	
	return srv.entitiesMap[id]
}

func (srv *Server) FindEntity(pred func (e IEntity) bool) IEntity {
	srv.entitiesLock.RLock()
	defer srv.entitiesLock.RUnlock()

	for _, e := range srv.entitiesMap {
		if (pred(e)) {
			return e
		}
	}
	
	return nil
}

func (srv *Server) GetEntityOwner(id uid.Uid) IOwnerEntity {
	srv.entitiesLock.RLock()
	defer srv.entitiesLock.RUnlock()
	
	entity := srv.entitiesMap[id]
	if entity == nil {
		return nil
	}

	ownerId := entity.GetOwner()
	if ownerId == uid.ServerUid {
		return nil
	}

	ownerEntity := srv.entitiesMap[ownerId]
	owner, isOwner := ownerEntity.(IOwnerEntity)
	if !isOwner {
		log.Warn().
			Str("entity_id", string(entity.GetId())).
			Any("owner_id", string(ownerId)).
			Msg("GetEntityOwner: entity has an invalid owner")
		return nil
	}

	return owner
}

func (srv *Server) ForEachEntity(cb func(player IEntity) (shouldStop bool)) {
	srv.entitiesLock.RLock()
	defer srv.entitiesLock.RUnlock()

	for _, player := range srv.entitiesMap {
		if cb(player) {
			break
		}
	}
}

func (srv *Server) AddReport(report model.IReport) {
	srv.reportsLock.Lock()
	defer srv.reportsLock.Unlock()

	if len(report.GetReport().ReportId) == 0 {
		panic("empty report id")
	}

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

// returns wether the given point is safe for player spawning
// also gives the average position of all graas tiles found
// (used as an heuristic for which direction to go after)
func (srv *Server) emptyProportion(point Point, freeAreaSide int) (bool, float64, float64) {
	count := 0.
	grass_count := 0
	sum_x := 0.
	sum_y := 0.

	for dx := -freeAreaSide; dx <= freeAreaSide; dx++ {
		for dy := -freeAreaSide; dy <= freeAreaSide; dy++ {
			np := point.Plus(dx, dy)

			tile := srv.tilemap.GetTile(np)
			if tile.Kind == terrain.TileGrass {
				sum_x += float64(np.X)
				sum_y += float64(np.Y)
				grass_count += 1
			}

			count++
		}
	}

	return int(count) == grass_count, sum_x / count, sum_y / count
}

// Returns a point safe for spawning near the given point and true
// or 0,0 and false it none could be found
//
// (yes kinda over engineered)
func (srv *Server) findSpawnPointNear(from Point, freeAreaSide int) (Point, bool) {
	visited := make(map[Point]bool)

	for {
		visited[from] = true

		all, a_x, a_y := srv.emptyProportion(from, freeAreaSide)
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

// Returns a safe spawn point with graas tile in the given cube "radius"
// also only consider a point if filter returns true
// reentry will deadlock but all other functions of server are available
func (srv *Server) FindSpawnPoint(start Point, freeAreaSide int, filter func(p Point) bool) Point {
	dist := 5

	srv.randLock.Lock()
	defer srv.randLock.Unlock()

	log.Trace().Int("dist", dist).Msg("[SPAWN_POINT] Looking for spawn point...")
	for dist < 1_000_000_000 {
		for try := 0; try < 120; try++ {
			center := start.Add(Point{
				X: srv.spawnRand.Intn(dist*2) - dist,
				Y: srv.spawnRand.Intn(dist*2) - dist,
			})
			point, found := srv.findSpawnPointNear(center, 2)

			if found && filter(point) {
				log.Trace().Any("point", point).Msg("[SPAWN_POINT] Found")
				return point
			}
		}
		log.Trace().Int("dist", dist).Msg("[SPAWN_POINT] not found, increasing dist")
		dist += dist / 2
	}

	panic("could not find spawn point")
}
