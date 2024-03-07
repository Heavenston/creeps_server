package main

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var CLI struct {
	Db string `short:"d" default:":memory:" help:"The sqlite connection string to the database"`

	Host string `short:"t" default:"localhost" help:"Target hostname for the api"`
	Port uint16 `short:"p" default:"16969" help:"Target port for the api"`

	Verbose int  `short:"v" type:"counter" help:"Once to enable debug logs, twice for trace logs"`
	Quiet   bool `short:"q" help:"If present overrides verbose and disables info logs and under"`
}

func main() {
	cw := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	log.Logger = zerolog.New(cw).With().
		Timestamp().
		Logger()

	ctx := kong.Parse(&CLI)
	if ctx.Error != nil {
		log.Fatal().Err(ctx.Error).Msg("CLI error")
	}

	db, err := gorm.Open(sqlite.Open(CLI.Db), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("DB Error")
	}

	db.AutoMigrate(&model.User{})

	log.Info().Str("url", CLI.Db).Msg("Connected to database")
}
