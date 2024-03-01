package events

import "sync/atomic"

// zero value is a valid not-cancelled handle
type CancelHandle struct {
	cancelled atomic.Bool
}

func (h *CancelHandle) Cancel() {
	h.cancelled.Store(true)
}

func (h *CancelHandle) IsCancelled() bool {
	return h.cancelled.Load()
}
