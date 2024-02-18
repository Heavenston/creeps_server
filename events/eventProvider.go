package events

import (
	"sync"
	"sync/atomic"
)

type CancelHandle struct {
    cancelled atomic.Bool
}

func (h *CancelHandle) Cancel() {
    h.cancelled.Store(true)
}

type sub[T any] struct {
    sendChan chan T
    handle *CancelHandle
}

// provider
// zero value is valid
type EventProvider[T any] struct {
    mutex sync.Mutex
    subs []sub[T]
}

func (provider *EventProvider[T]) Subscribe(channel chan T) *CancelHandle {
    provider.mutex.Lock()
    defer provider.mutex.Unlock()

    handle := new(CancelHandle)
    provider.subs = append(provider.subs, sub[T]{
        sendChan: channel,
        handle: handle,
    })
    handle.cancelled.Store(false)
    return handle
}

func (provider *EventProvider[T]) Emit(event T) {
    var i int = 0
    for i < len(provider.subs) {
        sub := &provider.subs[i]
        if sub.handle.cancelled.Load() {
            close(sub.sendChan)
            copy(provider.subs[i:], provider.subs[i+1:])
            provider.subs = provider.subs[:len(provider.subs)-1]
            continue
        }

        sub.sendChan <- event
        i++
    }
}
