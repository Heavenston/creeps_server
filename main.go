package main

import (
	"fmt"
	"time"

	"creeps.heav.fr/api"
	. "creeps.heav.fr/server"
	. "creeps.heav.fr/server/model"
	. "creeps.heav.fr/server/terrain"
	. "creeps.heav.fr/geom"
)

func main() {
	generator := NewChunkGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)
	srv := NewServer(&tilemap, &SetupResponse{
		TicksPerSeconds: 5,
	}, &CostsResponse{})

	api_server := &api.ApiServer {
		Addr: "localhost:1664",
		Server: srv,
	}

	srv.Ticker().AddTickFunc(func () {
		srv.Tilemap().PrintRegion(
			Point{X:-20,Y:-20},
			Point{X: 21,Y: 21},
		)
		fmt.Printf("\033[40A")
	})

	go api_server.Start()
	srv.Start()
}
