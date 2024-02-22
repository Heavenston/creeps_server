package entities

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
	"creeps.heav.fr/spatialmap"
	"creeps.heav.fr/uid"
	"github.com/rs/zerolog/log"
)

type Player struct {
	OwnerEntity

	server *Server

	// locks everything not read only
	lock sync.RWMutex

	id         uid.Uid
	username   string
	addr       string
	spawnPoint Point

	resources model.Resources

	townHalls []Point

	lastEnemySpawnTick int
}

func NewPlayer(
	server *Server,
	username string,
	addr string,
	spawnPoint Point,
) *Player {
	player := new(Player)

	player.OwnerEntity.InitOwnedEntities()

	player.server = server;
	player.spawnPoint = spawnPoint
	player.addr = addr
	player.id = uid.GenUid()
	player.username = username
	player.lastEnemySpawnTick = server.Ticker().GetTickNumber()

	return player
}

func (player *Player) GetServer() *Server {
	return player.server
}

func (player *Player) GetId() uid.Uid {
	return player.id
}

func (player *Player) GetUsername() string {
	return player.username
}

func (player *Player) GetAddr() string {
	return player.addr
}

func (player *Player) GetSpawnPoint() Point {
	return player.spawnPoint
}

// for IEntity
func (player *Player) GetAABB() AABB {
	return AABB{
		From: player.spawnPoint,
	}
}

// for IEntity
func (player *Player) GetOwner() uid.Uid {
	return uid.ServerUid
}

// for IEntity
func (player *Player) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return nil
}

func (player *Player) GetResources() model.Resources {
	player.lock.RLock()
	defer player.lock.RUnlock()
	return player.resources
}

// Do not call for modification after GetResources, to avoid race conditions use
// modify resources
func (player *Player) SetResources(resources model.Resources) {
	player.lock.Lock()
	defer player.lock.Unlock()

	player.resources = resources
}

// atomically modifies the resources
func (player *Player) ModifyResources(f func(res model.Resources) model.Resources) {
	player.lock.Lock()
	defer player.lock.Unlock()

	player.resources = f(player.resources)
}

func (player *Player) GetTownHalls() []Point {
	player.lock.RLock()
	defer player.lock.RUnlock()

	return player.townHalls
}

func (player *Player) hasTownHall(p Point) bool {
	for _, op := range player.townHalls {
		if op == p {
			return true
		}
	}
	return false
}

func (player *Player) HasTownHall(p Point) bool {
	player.lock.RLock()
	defer player.lock.RUnlock()
	return player.hasTownHall(p)
}

func (player *Player) AddTownHall(p Point) {
	player.lock.Lock()
	defer player.lock.Unlock()

	if player.hasTownHall(p) {
		log.Warn().
			Str("player_id", string(player.id)).
			Any("th_pos", p).
			Msg("Attempted to add a town hall twice to a player")
		return
	}

	player.townHalls = append(player.townHalls, p)
}

// returns true if it was removed
// false if it wasn't found
func (player *Player) RemoveTownHall(p Point) bool {
	player.lock.Lock()
	defer player.lock.Unlock()

	for i, op := range player.townHalls {
		if p == op {
			player.townHalls[i] = player.townHalls[len(player.townHalls)-1]
			player.townHalls = player.townHalls[:len(player.townHalls)-1]
			return true
		}
	}
	return false
}

func (player *Player) Register() {
	player.server.RegisterEntity(player)
	player.server.Events().Emit(&PlayerSpawnEvent{
		Player: player,
	})
}

func (player *Player) Unregister() {
	player.server.RemoveEntity(player.id)
	player.server.Events().Emit(&PlayerDespawnEvent{
		Player: player,
	})
}

// called each tick if enemy spawning is enabled
func (player *Player) enemySpawnTick() {
	currentTick := player.server.Ticker().GetTickNumber()

	elapsed := currentTick - player.lastEnemySpawnTick
	rate := player.server.GetSetup().EnemyBaseTickRate

	if elapsed <= rate {
		return
	}

	player.lastEnemySpawnTick = currentTick

	// finding spawn point can be costly so use another goroutine
	go (func () {
		NewRaid(player.server, player.id).Register()
	})()
}

func (player *Player) Tick() {
	player.lock.Lock()
	defer player.lock.Unlock()

	hasCitizens := false
	player.ForEachEntities(func(entity IEntity) (shouldStop bool) {
		if unit, ok := entity.(IUnit); ok {
			hasCitizens = unit.GetOpCode() == "citizen"
			shouldStop = hasCitizens
		}
		return
	})

	hasTownhalls := false
	for _, th := range player.townHalls {
		hasTownhalls = hasTownhalls ||
			player.server.Tilemap().GetTile(th).Kind == terrain.TileTownHall
	}

	if !hasCitizens || !hasTownhalls {
		player.Unregister();
		return
	}

	if player.server.GetSetup().EnableEnemies {
		player.enemySpawnTick()
	}
}
