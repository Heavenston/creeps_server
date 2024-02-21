package server

import (
	"time"

	"github.com/rs/zerolog/log"
)

type TickFunc func()

type Ticker struct {
	ticksPerSeconds float64

	tickNumber int
	startedAt  time.Time

	tickFuncs []TickFunc
}

func NewTicker(ticksPerSeconds float64) *Ticker {
	ticker := new(Ticker)
	ticker.startedAt = time.Now()
	ticker.tickNumber = 0
	ticker.tickFuncs = make([]TickFunc, 0)
	ticker.ticksPerSeconds = ticksPerSeconds
	return ticker
}

func (ticker *Ticker) Start() {
    log.Info().Float64("tps", ticker.ticksPerSeconds).Msg("Ticker starting")
	time_ticker := time.NewTicker(time.Duration(float64(time.Second) / ticker.ticksPerSeconds))
	defer time_ticker.Stop()

	for {
		for _, tf := range ticker.tickFuncs {
			go tf()
		}

		_ = <-time_ticker.C
		ticker.tickNumber++
	}
}

func (ticker *Ticker) GetTickNumber() int {
	return ticker.tickNumber
}

func (ticker *Ticker) AddTickFunc(f TickFunc) {
	ticker.tickFuncs = append(ticker.tickFuncs, f)
}