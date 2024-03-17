package events_test

import (
	"testing"

	"github.com/heavenston/creeps_server/creeps_lib/events"
)

func TestEventProvider(t *testing.T) {
    provider := events.EventProvider[int]{}
    provider.Emit(5)

    chan1 := make(chan int, 2)
    cancel := provider.Subscribe(chan1)

    select {
    case v := (<- chan1):
        t.Errorf("Received past event: %d", v)
        return
    default:
    }

    provider.Emit(10)
    provider.Emit(11)

    select {
    case v := (<- chan1):
        if v != 10 {
            t.Errorf("Received wrong value: %d", v)
        }
    default:
        t.Errorf("Did not receive value")
    }

    cancel.Cancel()
    provider.Emit(12)

    select {
    case v := (<- chan1):
        if v != 11 {
            t.Errorf("Received wrong value: %d", v)
        }
    default:
        t.Errorf("Did not receive value")
    }

    provider.Emit(13)

    select {
    case v, ok := (<- chan1):
        if !ok {
            break
        }
        t.Errorf("Should not have received value after cancel: %d", v)
    default:
    }
}
