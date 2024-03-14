package gamemanager

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
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
