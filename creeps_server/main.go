package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var tps float64
var apiPort int
var apiHost string
var viewerPort int
var viewerHost string

var rootCmd = cobra.Command {
	Use: "heav_creeps",
	Short: "A reimplementation of the *very* famous creeps game",

	Run: func(cmd *cobra.Command, args []string) {
		vl, _ := cmd.Flags().GetCount("verbose")
		log.Info().Int("vl", vl).Any("args", args).Msg("vl")
		if vl == 0 {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else if vl == 1 {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}
		startServ()
	},
}

func main() {
	cw := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	log.Logger = zerolog.New(cw).With().
		Timestamp().
		Logger()

	rootCmd.Flags().Float64VarP(&tps, "tps", "t", 5, "Ticks per seconds")
	rootCmd.Flags().IntVar(&apiPort, "api-port", 1664, "Port for the epita-compatible api")
	rootCmd.Flags().StringVar(&apiHost,
		"api-host", "",
		`Full host for the epita-compatible api, ex: 'localhost:1664'.
Overwrites api-port if present`,
	)
	rootCmd.Flags().IntVar(&viewerPort, "viewer-port", 1665, "Port for the viewer's api (not the viewer itself)")
	rootCmd.Flags().StringVar(&viewerHost,
		"viewer-host", "",
		`Full host for the viewer's api, ex: 'localhost:1665'.
Overwrites viewer-port if present`,
	)
	rootCmd.Flags().CountP("verbose", "v", "Once for debug twice for trace")
	rootCmd.Flags().BoolP("quiet", "q", false, "If present only warnings or errors are printed (overwrites verbose)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("cli error")
	}
}
