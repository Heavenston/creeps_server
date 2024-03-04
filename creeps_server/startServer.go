package main

import (
	"fmt"
	"time"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	. "github.com/heavenston/creeps_server/creeps_lib/terrain"
	"github.com/heavenston/creeps_server/creeps_server/epita_api"
	"github.com/heavenston/creeps_server/creeps_server/generator"
	. "github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/viewer"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func startServ() {
	switch logLevel(viper.GetString("loglevel")) {
	case logLevelTrace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case logLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case logLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case logLevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case logLevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	generator := generator.NewNoiseGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)

	var setup *model.SetupResponse
	viper.UnmarshalKey("setup", &setup)
	var costs *model.CostsResponse
	viper.UnmarshalKey("costs", &costs)

	if viper.IsSet("tps") {
		setup.TicksPerSecond = viper.GetFloat64("tps")
	}

	srv := NewServer(&tilemap, setup, costs)
	srv.SetDefaultPlayerResources(model.Resources{
		Rock:      30,
		Wood:      30,
		Food:      30,
		Oil:       0,
		Copper:    0,
		WoodPlank: 0,
	})

	api_server := &epita_api.ApiServer{
		Addr:   fmt.Sprintf("%s:%d", viper.GetString("api.host"), viper.GetInt("api.port")),
		Server: srv,
	}
	go api_server.Start()

	viewer_server := &viewer.ViewerServer{
		Addr:   fmt.Sprintf("%s:%d", viper.GetString("viewer.host"), viper.GetInt("viewer.port")),
		Server: srv,
	}
	go viewer_server.Start()

	tilemap.GenerateChunk(Point{X: 0, Y: 0})
	tilemap.GenerateChunk(Point{X: 0, Y: -1})
	tilemap.GenerateChunk(Point{X: -1, Y: 0})
	tilemap.GenerateChunk(Point{X: -1, Y: -1})

	srv.Start()
}
