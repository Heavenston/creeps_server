package main

import (
	"fmt"
	"time"

	"creeps.heav.fr/geom"
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	. "creeps.heav.fr/server/model"
	. "creeps.heav.fr/server/terrain"
	"creeps.heav.fr/units"
)

func main() {
	generator := NewChunkGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)
	srv := NewServer(&tilemap, &SetupResponse{
		TicksPerSeconds: 10,
	}, &CostsResponse{})

	player := NewPlayer("heavenstone")
	srv.RegisterPlayer(player)

	raider := units.NewRaiderUnit(srv, Point{X: 15, Y: 15})
	fmt.Printf("raider.GetId(): %v\n", raider.GetId())
	raider.SetPosition(geom.Point{X: 0, Y: 0})
	srv.RegisterUnit(raider)

	srv.Start()
}
