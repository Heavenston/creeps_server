package server

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/uid"
)

type Player struct {
	// locks everything not read only
	lock sync.RWMutex

	id        uid.Uid
	username  string
	addr      string
	spawnPoint Point

	resources model.Resources

	townHalls []Point
}

func NewPlayer(username string, addr string, spawnPoint Point) *Player {
	player := new(Player)

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
		panic("Cannot add a town hall twice")
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
