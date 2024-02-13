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
		EnableGC: false,
		GcTickRate: 50,
		EnableEnemies: false,
		EnemyTickRate: 75,
		MaxLoad: 25,
		MaxMissesPerPlayer: 200,
		MaxMissesPerUnit: 200,
		ServerId: "heavenstone_server",
		TicksPerSeconds: 5,
		TrackAchievements: false,
		WorldDimension: Point{},
		FoodGatherRate: 5,
		OilGatherRate: 5,
		RockGatherRate: 5,
		WoodGatherRate: 5,
	}, &model.CostsResponse{})
	srv.SetDefaultPlayerResources(model.Resources{
		Rock: 30,
		Wood: 30,
		Food: 30,
		Oil: 0,
		Copper: 0,
		WoodPlank: 0,
	})

	api_server := &api.ApiServer {
		Addr: "localhost:1664",
		Server: srv,
	}

	go api_server.Start()
	srv.Start()
}
