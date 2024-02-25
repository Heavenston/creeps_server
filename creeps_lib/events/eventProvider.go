package events

import (
	"sync"
)

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
    handle := new(CancelHandle)
    provider.SubscribeWithHandle(channel, handle)
    return handle
}

func (provider *EventProvider[T]) SubscribeWithHandle(channel chan T, handle *CancelHandle) {
    if handle.IsCancelled() {
        return;
    }

    provider.mutex.Lock()
    defer provider.mutex.Unlock()

    provider.subs = append(provider.subs, sub[T]{
        sendChan: channel,
        handle: handle,
    })
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
