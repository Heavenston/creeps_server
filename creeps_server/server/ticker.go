package server

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type TickFunc func()

type Ticker struct {
	ticksPerSecond float64

	tickNumber atomic.Int32
	startedAt  time.Time

	tickFuncsLock sync.RWMutex
	tickFuncs     []TickFunc

	deferedFuncsLock sync.RWMutex
	// see ticker.Defer
	deferedFuncs []TickFunc
}

func NewTicker(ticksPerSecond float64) *Ticker {
	ticker := new(Ticker)
	ticker.startedAt = time.Now()
	ticker.ticksPerSecond = ticksPerSecond
	return ticker
}

func (ticker *Ticker) Start() {
	log.Info().Float64("tps", ticker.ticksPerSecond).Msg("Ticker starting")
	time_ticker := time.NewTicker(ticker.TickDuration())
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
		ticker.tickNumber.Add(1)
	}
}

func (ticker *Ticker) GetTickNumber() int {
	return int(ticker.tickNumber.Load())
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

func (ticker *Ticker) TickDuration() time.Duration {
	return time.Duration(float64(time.Second) / ticker.ticksPerSecond)
}
