package server

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/uid"
	"github.com/rs/zerolog/log"
)

type Player struct {
	server *Server

	// locks everything not read only
	lock sync.RWMutex

	id         uid.Uid
	username   string
	addr       string
	spawnPoint Point

	resources model.Resources

	townHalls []Point
	units     []IUnit
}

func NewPlayer(
	server *Server,
	username string,
	addr string,
	spawnPoint Point,
) *Player {
	player := new(Player)

	player.server = server;
	player.spawnPoint = spawnPoint
	player.addr = addr
	player.id = uid.GenUid()
	player.username = username

	return player
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

func (player *Player) GetUnits() []IUnit {
	player.lock.RLock()
	defer player.lock.RUnlock()

	return player.units
}

func (player *Player) hasUnit(id uid.Uid) bool {
	for _, unit := range player.units {
		if unit.GetId() == id {
			return true
		}
	}
	return false
}

func (player *Player) HasUnit(id uid.Uid) bool {
	player.lock.RLock()
	defer player.lock.RUnlock()
	return player.hasUnit(id)
}

// this is done by the server on RegisterUnit
func (player *Player) AddUnit(p IUnit) {
	player.lock.Lock()
	defer player.lock.Unlock()

	if player.hasUnit(p.GetId()) {
		log.Warn().
			Str("player_id", string(player.id)).
			Str("unit_id", string(p.GetId())).
			Msg("Attempted to add a unit twice to a player")
		return
	}

	player.units = append(player.units, p)
}

// this is done by the server on RemoveUnit
func (player *Player) RemoveUnit(p IUnit) bool {
	player.lock.Lock()
	defer player.lock.Unlock()

	for i, unit := range player.units {
		if unit.GetId() == p.GetId() {
			player.units[i] = player.units[len(player.units)-1]
			player.units = player.units[:len(player.units)-1]
			return true
		}
	}
	return false
}

func (player *Player) kill() {
	player.server.RemovePlayer(player.id)
}

func (player *Player) Tick() {
	player.lock.Lock()
	defer player.lock.Unlock()

	if len(player.units) == 0 || len(player.townHalls) == 0 {
		player.kill();
		return
	}
}
