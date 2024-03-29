package main

import (
	"math"
	"os"

	"github.com/alecthomas/kong"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var CLI struct {
	ApiPort int16 `help:"Port for the epita-compatible api" default:"1664"`
	ApiHost string `help:"Host for the epita-compatible api" default:"localhost"`
	ViewerPort int16 `help:"Port for the viewer's api" default:"1665"`
	ViewerHost string `help:"Host for the viewer's api" default:"localhost"`

	Tps float64 `help:"Overrides the ticks per seconds"`
	Enemies *bool `negatable:"" help:"Overrides wether enemies are enables"`
	Hector *bool `negatable:"" help:"Overrides wether the garbage collector is enabled"`

	Verbose int `short:"v" type:"counter" help:"Once for debug prints, twice for trace"`
	Quiet bool `short:"q" help:"Overrites verbose, disables info logs and under"`
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
	cw := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	log.Logger = zerolog.New(cw).With().
		Timestamp().
		Logger()

	// viper.SetEnvPrefix("CREEPS")

	// viper.SetDefault("setup", &defaultSetup)
	// viper.SetDefault("costs", &defaultCosts)

	// rootCmd.Flags().Int("api-port", 1664, "")
	// viper.BindPFlag("api.port", rootCmd.Flags().Lookup("api-port"))
	// viper.BindEnv("api.port")

	// rootCmd.Flags().String(
	// 	"api-host", "localhost",
	// 	`Host (ip) for the epita-compitible api`,
	// )
	// viper.BindPFlag("api.host", rootCmd.Flags().Lookup("api-host"))
	// viper.BindEnv("api.host")

	// rootCmd.Flags().Int("viewer-port", 1665, "Port for the epita-compatible viewer")
	// viper.BindPFlag("viewer.port", rootCmd.Flags().Lookup("viewer-port"))
	// viper.BindEnv("viewer.port")

	// rootCmd.Flags().String(
	// 	"viewer-host", "localhost",
	// 	`Host (ip) for the viewer's api`,
	// )
	// viper.BindPFlag("viewer.host", rootCmd.Flags().Lookup("viewer-host"))
	// viper.BindEnv("viewer.host")

	// rootCmd.Flags().Float64("tps", defaultSetup.TicksPerSecond,
	// 	`Overwrites config's setup ticks per seconds`,
	// )
	// viper.BindPFlag("tps", rootCmd.Flags().Lookup("tps"))
	// viper.BindEnv("tps")
	
	// var level logLevel = logLevelInfo
	// rootCmd.Flags().VarP(&level, "loglevel", "l", `log level. Allowed values are "trace", "debug", "info", "warn" or "error"`)
	// viper.BindPFlag("loglevel", rootCmd.Flags().Lookup("loglevel"))
	// viper.BindEnv("loglevel", "LOGLEVEL")

	// viper.SetConfigName("heavcreeps") 
	// viper.SetConfigType("yaml") 
	// viper.AddConfigPath("$HOME/.heavcreeps")  
	// viper.AddConfigPath(".")               
	// err := viper.ReadInConfig() 
	// if err != nil {
	// 	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 	} else {
	// 		panic(fmt.Errorf("fatal error config file: %w", err))
	// 	}
	// }

	// if err := rootCmd.Execute(); err != nil {
	// 	log.Fatal().Err(err).Msg("cli error")
	// }

	ctx := kong.Parse(&CLI)
	startServ(ctx)
}
