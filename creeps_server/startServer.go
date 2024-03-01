package main

import (
	"fmt"
	"math"
	"time"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	. "github.com/heavenston/creeps_server/creeps_lib/terrain"
	"github.com/heavenston/creeps_server/creeps_server/epita_api"
	"github.com/heavenston/creeps_server/creeps_server/generator"
	. "github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/viewer"
)

func startServ() {
	generator := generator.NewNoiseGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)

	srv := NewServer(&tilemap, &model.SetupResponse{
		CitizenFeedingRate: 25,
		EnableGC:           false,
		GcTickRate:         150,
		EnableEnemies:      true,
		EnemyTickRate:      8,
		EnemyBaseTickRate:  300,
		MaxLoad:            20,
		MaxMissesPerPlayer: 200,
		MaxMissesPerUnit:   200,
		ServerId:           "heavenstone_server",
		TicksPerSeconds:    tps,
		TrackAchievements:  false,
		WorldDimension: Point{
			// big value but leave two bits to avoid any overflow anywhere
			X: math.MaxInt32 >> 2,
			Y: math.MaxInt32 >> 2,
		},
		FoodGatherRate: 5,
		OilGatherRate:  2,
		RockGatherRate: 5,
		WoodGatherRate: 5,
	}, &model.CostsResponse{
		BuildHousehold: model.CostResponse{
			Resources: model.Resources{
				Rock: 10,
				Wood: 10,
			},
			Cast: 6,
		},
		BuildRoad: model.CostResponse{
			Resources: model.Resources{
				Rock: 1,
			},
			Cast: 2,
		},
		BuildSawmill: model.CostResponse{
			Resources: model.Resources{
				Rock: 15,
				Wood: 25,
			},
			Cast: 10,
		},
		BuildSmeltery: model.CostResponse{
			Resources: model.Resources{
				Rock: 25,
				Wood: 15,
			},
			Cast: 2,
		},
		BuildTownHall: model.CostResponse{
			Resources: model.Resources{
				Rock: 100,
				Wood: 100,
			},
			Cast: 20,
		},
		Dismantle: model.CostResponse{
			Cast: 1,
		},
		Farm: model.CostResponse{
			Cast: 10,
		},
		FetchMessage: model.CostResponse{
			Cast: 1,
		},
		FireBomberBot: model.CostResponse{
			Cast: 6,
		},
		FireTurret: model.CostResponse{
			Cast: 2,
		},
		Gather: model.CostResponse{
			Cast: 4,
		},
		Move: model.CostResponse{
			Cast: 2,
		},
		Noop: model.CostResponse{
			Cast: 1,
		},
		Observe: model.CostResponse{
			Cast: 1,
		},
		RefineCopper: model.CostResponse{
			Resources: model.Resources{
				Rock: 10,
			},
			Cast: 8,
		},
		RefineWoodPlank: model.CostResponse{
			Resources: model.Resources{
				Wood: 10,
			},
			Cast: 8,
		},
		SendMessage: model.CostResponse{
			Cast: 1,
		},
		SpawnBomberBot: model.CostResponse{
			Resources: model.Resources{
				Rock: 5,
				Wood: 10,
			},
			Cast: 6,
		},
		SpawnTurret: model.CostResponse{
			Resources: model.Resources{
				Rock: 10,
				Wood: 5,
			},
			Cast: 6,
		},
		Unload: model.CostResponse{
			Cast: 1,
		},
		UpgradeBomberBot: model.CostResponse{
			Resources: model.Resources{
				Rock:      5,
				Wood:      10,
				Oil:       4,
				Copper:    1,
				WoodPlank: 2,
			},
			Cast: 1,
		},
		UpgradeCitizen: model.CostResponse{
			Resources: model.Resources{
				Rock:      5,
				Wood:      5,
				Food:      2,
				Copper:    1,
				WoodPlank: 1,
			},
			Cast: 1,
		},
		UpgradeTurret: model.CostResponse{
			Resources: model.Resources{
				Rock:      10,
				Wood:      5,
				Oil:       4,
				Copper:    3,
				WoodPlank: 1,
			},
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

	if apiHost == "" {
		apiHost = fmt.Sprintf("localhost:%d", apiPort)
	}

	api_server := &epita_api.ApiServer{
		Addr:   apiHost,
		Server: srv,
	}
	go api_server.Start()

	if viewerHost == "" {
		viewerHost = fmt.Sprintf("localhost:%d", viewerPort)
	}

	viewer_server := &viewer.ViewerServer{
		Addr:   viewerHost,
		Server: srv,
	}
	go viewer_server.Start()

	tilemap.GenerateChunk(Point{X: 0, Y: 0})
	tilemap.GenerateChunk(Point{X: 0, Y: -1})
	tilemap.GenerateChunk(Point{X: -1, Y: 0})
	tilemap.GenerateChunk(Point{X: -1, Y: -1})

	srv.Start()
}
