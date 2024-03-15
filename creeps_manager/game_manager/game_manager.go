package gamemanager

import (
	"fmt"
	"sync"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type GameManager struct {
	db         *gorm.DB
	binaryPath string

	gamesLock sync.Mutex
	games     map[int]*RunningGame

	port int
}

func NewGameManager(db *gorm.DB, binaryPath string) *GameManager {
	m := &GameManager{
		db:         db,
		binaryPath: binaryPath,

		games: make(map[int]*RunningGame),

		port: 1664,
	}
	return m
}

// Called once to go through games in the db
func (self *GameManager) Restore() error {
	time := time.Now()
	rs := self.db.Model(model.Game{}).
		Where("started_at IS NOT null AND ended_at IS null").
		Update("ended_at", &time)
	if rs.Error != nil {
		return rs.Error
	}

	log.Info().Any("count", rs.RowsAffected).Msg("Resored games")

	return nil
}

func (self *GameManager) GetRunningGame(id uint) *RunningGame {
	self.gamesLock.Lock()
	defer self.gamesLock.Unlock()
	return self.games[int(id)]
}

func (self *GameManager) forgetGame(rg *RunningGame) bool {
	self.gamesLock.Lock()
	defer self.gamesLock.Unlock()
	if self.games[rg.Id] != rg {
		return false
	}
	delete(self.games, rg.Id)
	return true
}

func (self *GameManager) StartGame(game model.Game) (*RunningGame, error) {
	self.gamesLock.Lock()
	defer self.gamesLock.Unlock()

	if self.games[int(game.ID)] != nil {
		return nil, fmt.Errorf("Game already started")
	}

	rg, err := newRunningGame(self, game, gameStartCfg{
		game: game,

		binaryPath: self.binaryPath,

		apiPort: self.port,
		viewerPort: self.port+1,
	})
	if err != nil {
		return nil, err
	}

	self.games[int(game.ID)] = rg
	self.port++

	return rg, nil
}
