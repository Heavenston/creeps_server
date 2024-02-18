package spatialmap

import (
	"sync"

	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
)

type MovedEvent struct {
	From Point
	To   Point
}

type Positioned interface {
	comparable
	MovementEvents() *events.EventProvider[MovedEvent]
	GetPosition() Point
}

type el[T Positioned] struct {
	val       T
	subHandle *events.CancelHandle
}

type internalMovedEvent[T Positioned] struct {
	From   Point
	To     Point
	Object T
}

// Operations for fetching a list of entities by their positions
type SpatialMap[T Positioned] struct {
	lock sync.RWMutex

	isCopy     bool
	updateChan chan internalMovedEvent[T]
	objects    []el[T]
}

func (m *SpatialMap[T]) updateRoutine() {
	for {
		_, ok := (<- m.updateChan)
		if !ok {
			break
		}
	}
}

func Make[T Positioned]() *SpatialMap[T] {
	this := &SpatialMap[T]{
		updateChan: make(chan internalMovedEvent[T]),
	}
	go this.updateRoutine()
	return this
}

// makes a shallow copy, the copy is read-only
func (m *SpatialMap[T]) Copy() SpatialMap[T] {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return SpatialMap[T]{
		isCopy: true,
		objects: m.objects,
	}
}

func (m *SpatialMap[T]) Add(p T) {
	if m.isCopy {
		panic("cannot modify a spatialmap copy")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, o := range m.objects {
		if o.val == p {
			panic("cannot have duplicate objects")
		}
	}
	channel := make(chan MovedEvent)
	handle := p.MovementEvents().Subscribe(channel)

	// 'converts' MovedEvents to internalMovedEvent
	go (func() {
		for {
			event, ok := (<-channel)
			if !ok {
				break
			}
			m.updateChan <- internalMovedEvent[T]{
				From:   event.From,
				To:     event.To,
				Object: p,
			}
		}
	})()

	m.objects = append(m.objects, el[T]{
		val:       p,
		subHandle: handle,
	})
}

func (m *SpatialMap[T]) RemoveFirst(predicate func(T) bool) *T {
	if m.isCopy {
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
