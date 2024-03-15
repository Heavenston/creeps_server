package gamemanager

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/rs/zerolog/log"
)

type RunningGame struct {
	Gm *GameManager

	Id        int
	CreatorId int

	ApiPort    int
	ViewerPort int

	Cmd      *exec.Cmd
}

type gameStartCfg struct {
	game model.Game

	binaryPath string

	apiPort int
	viewerPort int
}

func newRunningGame(gm *GameManager, game model.Game, cfg gameStartCfg) (*RunningGame, error) {
	logFile := fmt.Sprintf("/tmp/game%d.logs", game.ID)

	cmd := exec.Command(
		cfg.binaryPath,
		"--api-port", fmt.Sprintf("%d", cfg.apiPort),
		"--viewer-port", fmt.Sprintf("%d", cfg.viewerPort),
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
		Gm: gm,

		Id:        int(game.ID),
		CreatorId: game.CreatorID,

		ApiPort:    cfg.apiPort,
		ViewerPort: cfg.viewerPort,

		Cmd: cmd,
	}

	now := time.Now()
	gm.db.Model(&game).Where("id = ?", game.ID).Update("started_at", &now)

	log.Info().
		Str("binary_path", cfg.binaryPath).
		Str("logs", logFile).
		Str("name", game.Name).
		Uint("id", game.ID).
		Msg("Started game")

	return rgame, nil
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
