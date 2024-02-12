package main

import (
	"fmt"
	"time"

	. "creeps.heav.fr/server"
	. "creeps.heav.fr/server/model"
	. "creeps.heav.fr/server/terrain"
)

func main() {
	generator := NewChunkGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)
	srv := NewServer(&tilemap, &SetupResponse{
		TicksPerSeconds: 10,
	}, &CostsResponse{})

	// player := NewPlayer("heavenstone")
	// srv.RegisterPlayer(player)

	// raider := units.NewRaiderUnit(srv, Point{X: 15, Y: 15})
	// fmt.Printf("raider.GetId(): %v\n", raider.GetId())
	// raider.SetPosition(geom.Point{X: 0, Y: 0})
	// srv.RegisterUnit(raider)

	fmt.Println("looking for spawn point...")
	spawn := srv.FindSpawnPoint()
	fmt.Printf("found on %v\n", spawn)
	srv.Tilemap().PrintRegion(
		spawn.Plus(-20,-20),
		spawn.Plus( 21, 21),
	)

	// srv.Start()
}
