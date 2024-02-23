package main

import (
	"math"
	"os"
	"time"

	"creeps.heav.fr/epita_api"
	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	. "creeps.heav.fr/server/terrain"
	"creeps.heav.fr/viewer"
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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	generator := NewChunkGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)
	srv := NewServer(&tilemap, &model.SetupResponse{
		CitizenFeedingRate: 20,
		EnableGC:           false,
		GcTickRate:         50,
		EnableEnemies:      true,
		EnemyTickRate:      3,
		EnemyBaseTickRate:  150,
		MaxLoad:            40,
		MaxMissesPerPlayer: 200,
		MaxMissesPerUnit:   200,
		ServerId:           "heavenstone_server",
		TicksPerSeconds:    5,
		TrackAchievements:  false,
		WorldDimension:     Point{
			// big value but leave two bits to avoid any overflow anywhere
			X: math.MaxInt32 >> 2,
			Y: math.MaxInt32 >> 2,
		},
		FoodGatherRate:     5,
		OilGatherRate:      5,
		RockGatherRate:     5,
		WoodGatherRate:     5,
	}, &model.CostsResponse{
		BuildHousehold: model.CostResponse{
			Resources: model.Resources{
				Rock: 3,
				Wood: 5,
			},
			Cast: 1,
		},
		BuildRoad: model.CostResponse{
			Resources: model.Resources{
				Rock: 1,
			},
			Cast: 1,
		},
		BuildSawmill: model.CostResponse{
			Cast: 1,
		},
		BuildSmeltery: model.CostResponse{
			Cast: 1,
		},
		BuildTownHall: model.CostResponse{
			Cast: 1,
		},
		Dismantle: model.CostResponse{
			Cast: 1,
		},
		Farm: model.CostResponse{
			Cast: 1,
		},
		FetchMessage: model.CostResponse{
			Cast: 1,
		},
		FireBomberBot: model.CostResponse{
			Cast: 1,
		},
		FireTurret: model.CostResponse{
			Cast: 1,
		},
		Gather: model.CostResponse{
			Cast: 1,
		},
		Move: model.CostResponse{
			Cast: 1,
		},
		Noop: model.CostResponse{
			Cast: 1,
		},
		Observe: model.CostResponse{
			Cast: 1,
		},
		RefineCopper: model.CostResponse{
			Cast: 1,
		},
		RefineWoodPlank: model.CostResponse{
			Cast: 1,
		},
		SendMessage: model.CostResponse{
			Cast: 1,
		},
		SpawnBomberBot: model.CostResponse{
			Cast: 1,
		},
		SpawnTurret: model.CostResponse{
			Cast: 1,
		},
		Unload: model.CostResponse{
			Cast: 1,
		},
		UpgradeBomberBot: model.CostResponse{
			Cast: 1,
		},
		UpgradeCitizen: model.CostResponse{
			Cast: 1,
		},
		UpgradeTurret: model.CostResponse{
			Cast: 1,
		},
	})
	srv.SetDefaultPlayerResources(model.Resources{
		Rock:      30,
		Wood:      30,
		Food:      30,
		Oil:       0,
		Copper:    0,
		WoodPlank: 0,
	})

	api_server := &epita_api.ApiServer{
		Addr:   "localhost:1664",
		Server: srv,
	}
	go api_server.Start()

	viewer_server := &viewer.ViewerServer{
		Addr:   "localhost:1665",
		Server: srv,
	}
	go viewer_server.Start()

	srv.Start()
}
