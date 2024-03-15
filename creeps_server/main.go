package main

import (
	"math"

	"github.com/alecthomas/kong"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
)

var CLI struct {
	ApiPort        int16   `help:"Port for the epita-compatible api" default:"0"`
	ApiHost        string  `help:"Host for the epita-compatible api" default:"localhost"`
	ApiPortFile    *string `help:"If given, the resolved port of the epita-compatible api will be written in it"`
	ViewerPort     int16   `help:"Port for the viewer's api" default:"0"`
	ViewerHost     string  `help:"Host for the viewer's api" default:"localhost"`
	ViewerPortFile *string `help:"If given, the revolved port of the viewer api will be written in it"`

	Tps     float64 `help:"Overrides the ticks per seconds"`
	Enemies *bool   `negatable:"" help:"Overrides wether enemies are enables"`
	Hector  *bool   `negatable:"" help:"Overrides wether the garbage collector is enabled"`

	LogFile *string `short:"l" help:"File in which to print logs instead of stdout"`

	Verbose int  `short:"v" type:"counter" help:"Once for debug prints, twice for trace"`
	Quiet   bool `short:"q" help:"Overrites verbose, disables info logs and under"`
}

var defaultSetup model.SetupResponse = model.SetupResponse{
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
	TicksPerSecond:     5,
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
}

var defaultCosts model.CostsResponse = model.CostsResponse{
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
}

var defaultPlayerResources model.Resources = model.Resources{
	Rock:      30,
	Wood:      30,
	Food:      30,
	Oil:       0,
	Copper:    0,
	WoodPlank: 0,
}

func main() {
	ctx := kong.Parse(&CLI)
	startServ(ctx)
}
