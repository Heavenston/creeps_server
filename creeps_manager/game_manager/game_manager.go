package gamemanager

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RunningGame struct {
	Gm *GameManager

	Id        int
	CreatorId int

	ApiPort    int
	ViewerPort int

	Cmd      *exec.Cmd
}

func (self *RunningGame) Stop() error {
	if !self.Gm.forgetGame(self) {
		return fmt.Errorf("Running game already forgotten by game manager")
	}
	
	if self.Cmd.Process == nil {
		return fmt.Errorf("Process not started")
	}

	err := syscall.Kill(self.Cmd.Process.Pid, syscall.SIGTERM)
	if err != nil {
		return err
	}

	self.Gm.db.
		Model(&model.Game{}).
		Where("id = ?", self.Id).
		Update("ended_at", time.Now())

	return nil
}

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

	apiPort := self.port
	viewerPort := self.port + 1
	self.port += 2

	logFile := fmt.Sprintf("/tmp/game%d.logs", game.ID)

	cmd := exec.Command(
		self.binaryPath,
		"--api-port", fmt.Sprintf("%d", apiPort),
		"--viewer-port", fmt.Sprintf("%d", viewerPort),
		"-vv",
		"--log-file", logFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	rgame := &RunningGame{
		Gm: self,

		Id:        int(game.ID),
		CreatorId: game.CreatorID,

		ApiPort:    apiPort,
		ViewerPort: viewerPort,

		Cmd: cmd,
	}

	self.games[int(game.ID)] = rgame

	now := time.Now()
	self.db.Model(&game).Where("id = ?", game.ID).Update("started_at", &now)

	log.Info().
		Str("binary_path", self.binaryPath).
		Str("logs", logFile).
		Str("name", game.Name).
		Uint("id", game.ID).
		Msg("Started game")

	return rgame, nil
}
