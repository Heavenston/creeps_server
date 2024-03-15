package gamemanager

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
	viewerapimodel "github.com/heavenston/creeps_server/creeps_lib/viewer_api_model"
	"github.com/rs/zerolog/log"
)

type RunningGame struct {
	Gm *GameManager

	Id        int
	CreatorId int

	Host string
	ApiPort    int
	ViewerPort int

	Packets chan viewerapimodel.Message

	Cmd      *exec.Cmd
}

type gameStartCfg struct {
	game model.Game

	binaryPath string
}

func waitAndReadPortFile(filePath string) (uint16, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	file.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return 0, err
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		return 0, err
	}

	for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                break
            }
			if event.Has(fsnotify.Write) {
				break
			}
			continue
        case err, ok := <-watcher.Errors:
            if !ok {
                break
            }
			return 0, err
        }
		break
    }

	var port uint16
	file, err = os.Open(filePath)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	_, err = fmt.Fscanf(file, "%d", &port)
	if err != nil {
		return 0, err
	}

	return port, nil
}

func newRunningGame(gm *GameManager, game model.Game, cfg gameStartCfg) (*RunningGame, error) {
	logFile := fmt.Sprintf("/tmp/game%d.logs", game.ID)

	viewerPortFile := fmt.Sprintf("/tmp/game%d_viewer_port_%s", game.ID, uid.GenUid())
	apiPortFile := fmt.Sprintf("/tmp/game%d_api_port_%s", game.ID, uid.GenUid())

	viewerPort := make(chan uint16)
	go func() {
		port, err := waitAndReadPortFile(viewerPortFile)
		if err != nil {
			log.Error().Err(err).Send()
		}
		viewerPort <- port
	}()

	apiPort := make(chan uint16)
	go func() {
		port, err := waitAndReadPortFile(apiPortFile)
		if err != nil {
			log.Error().Err(err).Send()
		}
		apiPort <- port
	}()

	passwordAlphabet := "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ0123456789"
	passwordBytes := make([]byte, 100)
	_, err := rand.Read(passwordBytes)
	if err != nil {
		return nil, err
	}
	passBuilder := strings.Builder{}
	for _, b := range passwordBytes {
		passBuilder.WriteByte(passwordAlphabet[int(b) % len(passwordAlphabet)])
	}
	password := passBuilder.String()

	cmd := exec.Command(
		cfg.binaryPath,
		"--viewer-port-file", viewerPortFile,
		"--api-port-file", apiPortFile,
		"--viewer-admin-password", password,
		"--viewer-admin-host", "127.0.0.1",
		"-vv",
		"--log-file", logFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	rgame := &RunningGame{
		Gm: gm,

		Id:        int(game.ID),
		CreatorId: game.CreatorID,

		Host: "localhost",
		ViewerPort: int(<-viewerPort),
		ApiPort: int(<-apiPort),

		Packets: make(chan viewerapimodel.Message),

		Cmd: cmd,
	}

	err = rgame.connect(password)
	if err != nil {
		cmd.Process.Kill()
		return nil, err
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

// called by newRunningGame
func (self *RunningGame) connect(password string) error {
	addr := fmt.Sprintf(
		"ws://%s/websocket",
		net.JoinHostPort(self.Host, fmt.Sprint(self.ViewerPort)),
	)

	log.Debug().Str("ws", addr).Msg("Connecting to game server ws")
	
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", addr, err)
	}

	err = conn.WriteJSON(viewerapimodel.CreateMessage(viewerapimodel.C2SInit {
		AuthPassword: password,
	}))
	if err != nil {
		return fmt.Errorf("Init write error: %w", err)
	}

	var msg viewerapimodel.Message
	err = conn.ReadJSON(&msg)
	if err != nil {
		return fmt.Errorf("Init read error: %w", err)
	}

	if msg.Kind != (viewerapimodel.S2CInit{}.MsgKind()) {
		return fmt.Errorf("Received invalid handshake")
	}

	var init viewerapimodel.S2CInit
	err = json.Unmarshal(msg.Content, &init)
	if err != nil {
		return fmt.Errorf("Invalid init content: %w", err)
	}

	log.Debug().
		Int("game_id", self.Id).
		Str("addr", addr).
		Msg("Fully connected to running game")

	go func() {
		defer conn.Close()
		for {
			var msg viewerapimodel.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Warn().Err(err).Msg("Game connection error")
				return
			}

			log.Info().Any("msg", msg).Send()
		}
	}()

	return nil
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
