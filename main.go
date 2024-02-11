package main

import (
	"time"

	"creeps.heav.fr/server"
	"creeps.heav.fr/server/model"
	"creeps.heav.fr/server/terrain"
)

func main() {
	generator := terrain.NewChunkGenerator(time.Now().UnixMilli())
	tilemap := terrain.NewTilemap(generator)
	server := server.NewServer(&tilemap, &model.SetupResponse{
		TicksPerSeconds: 10,
	}, &model.CostsResponse{
		
	})
	server.Start()
}
