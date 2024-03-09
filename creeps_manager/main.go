package main

import (
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Heavenston/creeps_server/creeps_manager/discordapi"
	gamemanager "github.com/Heavenston/creeps_server/creeps_manager/game_manager"
	"github.com/Heavenston/creeps_server/creeps_manager/keys"
	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/Heavenston/creeps_server/creeps_manager/webserver"
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var CLI struct {
	Db string `env:"CREEPS_MANAGER_DB" short:"d" default:":memory:" help:"The sqlite connection string to the database"`

	Host string `short:"t" default:"localhost" help:"Target hostname for the api"`
	Port uint16 `short:"p" default:"1234" help:"Target port for the api"`

	ClientId         string `env:"CREEPS_MANAGER_CLIENT_ID" required:"" help:"Discord client id"`
	ClientSecret     string `env:"CREEPS_MANAGER_CLIENT_SECRET" required:"" help:"Discord client secret"`
	ServerBinaryPath string `env:"CREEPS_MANAGER_SERVER_BINARY" default:"creeps_server" help:"Path to the binary of the creeps server"`
	LoginURL *url.URL `env:"LOGIN_URL" required:"" help:"Discord Login url"`

	JWTSecret *string `env:"CREEPS_MANAGER_JWT_SECRET" help:"If present this overrides the generated jwt secret, intended for debugging use only"`

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

	if CLI.JWTSecret != nil {
		keys.SetJwtSecret([]byte(*CLI.JWTSecret))
	}

	switch CLI.Verbose {
	case 0:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	db, err := gorm.Open(sqlite.Open(CLI.Db), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("DB Error")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Role{})
	db.AutoMigrate(&model.Game{})

	log.Info().Str("url", CLI.Db).Msg("Connected to database")

	gameManager := gamemanager.NewGameManager(db, CLI.ServerBinaryPath)

	err = (&webserver.WebServer{
		Db:          db,
		GameManager: gameManager,
		LoginURL: CLI.LoginURL,

		DiscordAuth: &discordapi.DiscordAppAuth{
			ClientId:     CLI.ClientId,
			ClientSecret: CLI.ClientSecret,
		},
	}).Start(fmt.Sprintf("%s:%d", CLI.Host, CLI.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("HTTP Start Error")
	}
}
