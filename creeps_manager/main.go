package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = cobra.Command{
	Use:   "heav_creeps",
	Short: "A reimplementation of the *very* famous creeps game",

	Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("%s\n", cmd.Flag("executable").Value)
        fmt.Printf("%s\n", viper.GetString("executable"))
	},
}

func main() {
	cw := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	log.Logger = zerolog.New(cw).With().
		Timestamp().
		Logger()

    viper.SetEnvPrefix("CREEPS_MANAGER")

    rootCmd.Flags().StringP("executable", "e", "", "Path to the creeps game server executable")
    rootCmd.MarkFlagRequired("executable")
    viper.BindPFlag("executable", rootCmd.Flag("executable"))
    viper.BindEnv("executable")

    rootCmd.Flags().Int16P("port", "p", 1777, "Target port")
    viper.BindEnv("port")
    rootCmd.Flags().StringP("host", "t", "localhost", "Target host")
    viper.BindEnv("host")

	viper.SetConfigName("heavcreepsmanager") 
	viper.SetConfigType("yaml") 
	viper.AddConfigPath("$HOME/.heavcreepsmanager")  
	viper.AddConfigPath(".")               
	err := viper.ReadInConfig() 
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("cli parse error")
	}
}
