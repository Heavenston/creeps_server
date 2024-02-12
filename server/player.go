package server

import (
    "creeps.heav.fr/server/model"
    . "creeps.heav.fr/geom"
)

type Player struct {
    id Uid
    username string
    resources model.Resources

    townHalls []Point
}

func NewPlayer(username string) *Player {
    player := new(Player)

    player.id = GenUid()
    player.username = username

    return player
}

func (player *Player) GetId() Uid {
    return player.id
}

func (player *Player) GetUsername() string {
    return player.username
}

func (player *Player) GetResources() *model.Resources {
    return &player.resources
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
