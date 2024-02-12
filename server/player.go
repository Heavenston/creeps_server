package server

import (
	"creeps.heav.fr/uid"
	. "creeps.heav.fr/geom"
	"creeps.heav.fr/api/model"
)

type Player struct {
	id        uid.Uid
	username  string
	resources model.Resources
	addr      string

	townHalls []Point
}

func NewPlayer(username string, addr string) *Player {
	player := new(Player)

	player.addr = addr
	player.id = uid.GenUid()
	player.username = username

	return player
}

func (player *Player) GetId() uid.Uid {
	return player.id
}

func (player *Player) GetAddr() string {
	return player.addr
}

func (player *Player) GetUsername() string {
	return player.username
}

func (player *Player) GetResources() model.Resources {
	return player.resources
}

func (player *Player) SetResources(resources model.Resources)  {
	player.resources = resources
}

func (player *Player) GetTownHalls() []Point {
	return player.townHalls
}

func (player *Player) HasTownHall(p Point) bool {
	for _, op := range player.townHalls {
		if op == p {
			return true
		}
	}
	return false
}

func (player *Player) AddTownHall(p Point) {
	if player.HasTownHall(p) {
		panic("Cannot add a town hall twice")
	}

	player.townHalls = append(player.townHalls, p)
}

// returns true if it was removed
// false if it wasn't found
func (player *Player) RemoveTownHall(p Point) bool {
	for i, op := range player.townHalls {
		if p == op {
			player.townHalls[i] = player.townHalls[len(player.townHalls)-1]
			player.townHalls = player.townHalls[:len(player.townHalls)-1]
			return true
		}
	}
	return false
}
