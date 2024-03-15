package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alecthomas/kong"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	. "github.com/heavenston/creeps_server/creeps_lib/terrain"
	"github.com/heavenston/creeps_server/creeps_server/epita_api"
	"github.com/heavenston/creeps_server/creeps_server/generator"
	. "github.com/heavenston/creeps_server/creeps_server/server"
	"github.com/heavenston/creeps_server/creeps_server/viewer_api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func startServ(*kong.Context) {
	var cw io.Writer
	if CLI.LogFile == nil {
		cw = zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
	} else {
		var err error
		cw, err = os.Create(*CLI.LogFile)
		if err != nil {
			panic(err)
		}
	}
	log.Logger = zerolog.New(cw).With().
		Caller().
		Timestamp().
		Logger()

	switch CLI.Verbose {
	case 0:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	if CLI.Quiet {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	generator := generator.NewNoiseGenerator(time.Now().UnixMilli())
	tilemap := NewTilemap(generator)

	setup := defaultSetup
	costs := defaultCosts

	if CLI.Tps > 0 {
		setup.TicksPerSecond = CLI.Tps
	}
	if CLI.Hector != nil {
		setup.EnableGC = *CLI.Hector
	}
	if CLI.Enemies != nil {
		setup.EnableEnemies = *CLI.Enemies
	}

	srv := NewServer(&tilemap, &setup, &costs)
	srv.SetDefaultPlayerResources(model.Resources{
		Rock:      30,
		Wood:      30,
		Food:      30,
		Oil:       0,
		Copper:    0,
		WoodPlank: 0,
	})

	api_server := &epita_api.ApiServer{ Server: srv }
	go func() {
		err := api_server.Start(
			fmt.Sprintf("%s:%d", CLI.ApiHost, CLI.ApiPort),
			CLI.ApiPortFile,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("could not start api server")
		}
	}()

	viewer_server := &viewer_api.ViewerServer{ Server: srv }
	go func() {
		err := viewer_server.Start(
			fmt.Sprintf("%s:%d", CLI.ViewerHost, CLI.ViewerPort),
			CLI.ViewerPortFile,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("could not start viewer server")
		}
	}()

	tilemap.GenerateChunk(Point{X: 0, Y: 0})
	tilemap.GenerateChunk(Point{X: 0, Y: -1})
	tilemap.GenerateChunk(Point{X: -1, Y: 0})
	tilemap.GenerateChunk(Point{X: -1, Y: -1})

	srv.Start()
}
