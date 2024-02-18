package spatialmap

import (
	"sync"

	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
)

type ObjectMovedEvent struct {
	From Point
	To   Point
}

type Positioned interface {
	comparable
	MovementEvents() *events.EventProvider[ObjectMovedEvent]
	GetPosition() Point
}

type el[T Positioned] struct {
	val       T
	subHandle *events.CancelHandle
}

type SpatialMapEvent[T Positioned] struct {
	PreviousPosition Point
	Position         Point
	Object           T
}

type sub[T Positioned] struct {
	From    Point
	Upto    Point
	Channel chan SpatialMapEvent[T]
	Handle  *events.CancelHandle
}

// Data structure for fetching a list of entities by their positions
// TODO: Optimize with like maybe an octree or smthg
type SpatialMap[T Positioned] struct {
	lock sync.RWMutex

	isReadOnly bool
	updateChan chan SpatialMapEvent[T]
	objects    []el[T]

	subscriptions []sub[T]
}

// rountine always running while the spatial map is active
func (m *SpatialMap[T]) updateRoutine() {
	for {
		event, ok := (<-m.updateChan)
		if !ok {
			break
		}

		// also filters the subscriptions
		index := 0
		for _, sub := range m.subscriptions {
			// filter
			if !sub.Handle.IsCancelled() {
				m.subscriptions[index] = sub
				index++
			} else {
				close(sub.Channel)
				continue
			}

			if !event.Position.IsWithing(sub.From, sub.Upto) && !event.PreviousPosition.IsWithing(sub.From, sub.Upto) {
				continue
			}
			sub.Channel <- event
		}
		m.subscriptions = m.subscriptions[:index]
	}
}

// must be eventually closed with Close, or a goroutine leak will occure
func Make[T Positioned]() *SpatialMap[T] {
	this := &SpatialMap[T]{
		updateChan: make(chan SpatialMapEvent[T]),
	}
	go this.updateRoutine()
	return this
}

// after closing the map is now readonly
func (m *SpatialMap[T]) Close() {
	close(m.updateChan)
	m.isReadOnly = true
}

// makes a shallow copy, the copy is read-only
func (m *SpatialMap[T]) Copy() SpatialMap[T] {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return SpatialMap[T]{
		isReadOnly: true,
		objects:    m.objects,
	}
}

func (m *SpatialMap[T]) Add(p T) {
	if m.isReadOnly {
		panic("cannot modify a spatialmap copy")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, o := range m.objects {
		if o.val == p {
			panic("cannot have duplicate objects")
		}
	}
	channel := make(chan ObjectMovedEvent)
	handle := p.MovementEvents().Subscribe(channel)

	// 'converts' MovedEvents to internalMovedEvent
	go (func() {
		for {
			event, ok := (<-channel)
			if !ok {
				break
			}
			m.updateChan <- SpatialMapEvent[T]{
				PreviousPosition: event.From,
				Position:         event.To,
				Object:           p,
			}
		}
	})()

	m.objects = append(m.objects, el[T]{
		val:       p,
		subHandle: handle,
	})
}

func (m *SpatialMap[T]) RemoveFirst(predicate func(T) bool) *T {
	if m.isReadOnly {
		panic("cannot modify a spatialmap copy")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	for i, o := range m.objects {
		if predicate(o.val) {
			o.subHandle.Cancel()
			// swap with last remove
			m.objects[i] = m.objects[len(m.objects)-1]
			m.objects = m.objects[:len(m.objects)-1]
			return &o.val
		}
	}
	return nil
}

func (m *SpatialMap[T]) Find(predicate func(T) bool) *T {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, obj := range m.objects {
		if predicate(obj.val) {
			return &obj.val
		}
	}
	return nil
}

func (m *SpatialMap[T]) GetAt(point Point) *T {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.Find(func(t T) bool {
		return t.GetPosition() == point
	})
}

func (m *SpatialMap[T]) GetAllWithin(from Point, upto Point) []T {
	m.lock.RLock()
	defer m.lock.RUnlock()

	result := make([]T, 0)

	for _, obj := range m.objects {
		if obj.val.GetPosition().IsWithing(from, upto) {
			result = append(result, obj.val)
		}
	}

	return result
}

// subscribes to all objects moving from, into, or within the given region
//
// will send an event with the same value for PreviousPosition and Position with
// all objects that were already inside
func (m *SpatialMap[T]) SubscribeWithin(
	from Point,
	upto Point,
	channel chan SpatialMapEvent[T],
) *events.CancelHandle {
	m.lock.Lock()
	defer m.lock.Unlock()
	handle := &events.CancelHandle{}

	m.subscriptions = append(m.subscriptions, sub[T]{
		From:    from,
		Upto:    upto,
		Channel: channel,
		Handle:  handle,
	})

	for _, obj := range m.GetAllWithin(from, upto) {
		pos := obj.GetPosition()
		channel <- SpatialMapEvent[T] {
			PreviousPosition: pos,
			Position: pos,
			Object: obj,
		}
	}

	return handle
}

func (m *SpatialMap[T]) Iter() func() (bool, int, *T) {
	var i int = 0
	return (func() (bool, int, *T) {
		m.lock.RLock()
		defer m.lock.RUnlock()

		if i == len(m.objects) {
			return false, 0, nil
		}
		v := m.objects[i]
		i += 1
		return true, i - 1, &v.val
	})
}
