package main

import (
	"fmt"
	"time"

	"creeps.heav.fr/server"
	"creeps.heav.fr/server/model"
	"creeps.heav.fr/server/terrain"
)

func main() {
	generator := terrain.NewChunkGenerator(time.Now().UnixMilli())
	tilemap := terrain.NewTilemap(generator)
	srv := server.NewServer(&tilemap, &model.SetupResponse{
		TicksPerSeconds: 10,
	}, &model.CostsResponse{})

	raider := server.NewRaiderUnit()
	srv.RegisterUnit(raider)
	fmt.Printf("raider.GetId(): %v\n", raider.GetId())

	raider = server.NewRaiderUnit()
	srv.RegisterUnit(raider)
	fmt.Printf("raider.GetId(): %v\n", raider.GetId())

	srv.Start()
}
