package main

import (
	"os"
	"time"

	"creeps.heav.fr/api"
	"creeps.heav.fr/api/model"
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	. "creeps.heav.fr/server/terrain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cw := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	log.Logger = zerolog.New(cw).With().
		Timestamp().
		Logger()

	generator := NewChunkGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)
	srv := NewServer(&tilemap, &model.SetupResponse{
		CitizenFeedingRate: 20,
		EnableGC: true,
		GcTickRate: 120,
		EnableEnemies: false,
		EnemyTickRate: 0,
		FoodGatherRate: 0,
		MaxLoad: 50,
		MaxMissesPerPlayer: 200,
		MaxMissesPerUnit: 200,
		ServerId: "id",
		TicksPerSeconds: 5,
		TrackAchievements: false,
		WorldDimension: Point{},
		OilGatherRate: 5,
		RockGatherRate: 5,
		WoodGatherRate: 5,
	}, &model.CostsResponse{})

	api_server := &api.ApiServer {
		Addr: "localhost:1664",
		Server: srv,
	}

	// srv.Ticker().AddTickFunc(func () {
	// 	fmt.Printf("\033[40A")
	// 	srv.Tilemap().PrintRegion(
	// 		Point{X:-20,Y:-20},
	// 		Point{X: 21,Y: 21},
	// 	)
	// })

	go api_server.Start()
	srv.Start()
}
