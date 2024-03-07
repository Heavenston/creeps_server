package gamemanager

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"gorm.io/gorm"
)

type RunningGame struct {
	Id        int
	CreatorId int

	ApiPort    int
	ViewerPort int

	Cmd *exec.Cmd
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

func (self *GameManager) GetRunningGame(id uint) *RunningGame {
	self.gamesLock.Lock()
	defer self.gamesLock.Unlock()
	return self.games[int(id)]
}

func (self *GameManager) StartGame(game model.Game) (*RunningGame, error) {
	self.gamesLock.Lock()
	defer self.gamesLock.Unlock()

	apiPort := self.port
	viewerPort := self.port + 1
	self.port += 2

	cmd := exec.Command(
		self.binaryPath,
		"--api-port", fmt.Sprintf("%d", apiPort),
		"--viewer-port", fmt.Sprintf("%d", viewerPort),
	)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	rgame := &RunningGame{
		Id:        int(game.ID),
		CreatorId: game.CreatorID,

		ApiPort:    apiPort,
		ViewerPort: viewerPort,

		Cmd: cmd,
	}

	self.games[int(game.ID)] = rgame

	now := time.Now()
	self.db.Model(&game).Where("id = ?", game.ID).Update("started_at", &now)

	return rgame, nil
}
