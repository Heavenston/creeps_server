package server

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type TickFunc func()

type Ticker struct {
	ticksPerSeconds float64

	tickNumber int
	startedAt  time.Time

	tickFuncsLock sync.RWMutex
	tickFuncs     []TickFunc

	deferedFuncsLock sync.RWMutex
	// see ticker.Defer
	deferedFuncs []TickFunc
}

func NewTicker(ticksPerSeconds float64) *Ticker {
	ticker := new(Ticker)
	ticker.startedAt = time.Now()
	ticker.tickNumber = 0
	ticker.ticksPerSeconds = ticksPerSeconds
	return ticker
}

func (ticker *Ticker) Start() {
	log.Info().Float64("tps", ticker.ticksPerSeconds).Msg("Ticker starting")
	time_ticker := time.NewTicker(time.Duration(float64(time.Second) / ticker.ticksPerSeconds))
	defer time_ticker.Stop()


	for {
		start := time.Now()
		log.Trace().Msg("Started tick")

		tickFuncs := make([]TickFunc, len(ticker.tickFuncs))
		// copy to release the lock during the tick
		ticker.tickFuncsLock.RLock()
		copy(tickFuncs, ticker.tickFuncs)
		ticker.tickFuncsLock.RUnlock()

		for _, fun := range tickFuncs {
			fun()
		}

		ticker.deferedFuncsLock.Lock()
		defered := ticker.deferedFuncs
		ticker.deferedFuncs = nil
		ticker.deferedFuncsLock.Unlock()

		for _, fun := range defered {
			fun()
		}

		log.Trace().TimeDiff("took", time.Now(), start).Msg("Finished tick")

		_ = <-time_ticker.C
		ticker.tickNumber++
	}
}

func (ticker *Ticker) GetTickNumber() int {
	return ticker.tickNumber
}

func (ticker *Ticker) AddTickFunc(f TickFunc) {
	ticker.tickFuncsLock.Lock()
	defer ticker.tickFuncsLock.Unlock()
	ticker.tickFuncs = append(ticker.tickFuncs, f)
}

// schedule the given function to be called at the end of the tick
func (ticker *Ticker) Defer(f TickFunc) {
	ticker.deferedFuncsLock.Lock()
	defer ticker.deferedFuncsLock.Unlock()
	ticker.deferedFuncs = append(ticker.deferedFuncs, f)
}
