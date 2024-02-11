package main

import (
	"fmt"
	"time"

	. "creeps.heav.fr/geom"
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

	raider := NewRaiderUnit(srv, Point{X: 15, Y: 15})
	fmt.Printf("raider.GetId(): %v\n", raider.GetId())

	raider = NewRaiderUnit(srv, Point{X: 10, Y: 10})
	fmt.Printf("raider.GetId(): %v\n", raider.GetId())

	srv.Start()
}
