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

type Spatialized interface {
	comparable
	// can return nil if the object's aabb is guarenteed to never change
	MovementEvents() *events.EventProvider[ObjectMovedEvent]
	// if the returned AABB has size 0 this the object will match all queries
	GetAABB() AABB
}

type el[T Spatialized] struct {
	val       T
	subHandle *events.CancelHandle
}

type SpatialMapEvent[T Spatialized] struct {
	PreviousPosition Point
	Position         Point
	Object           T
}

// Data structure for fetching a list of entities by their aabbs
// TODO: Optimize with like maybe an quadtree or smthg
type SpatialMap[T Spatialized] struct {
	lock sync.RWMutex

	isReadOnly bool
	updateChan chan SpatialMapEvent[T]
	objects    []el[T]
}

// rountine always running while the spatial map is active
func (m *SpatialMap[T]) updateRoutine() {
	for {
		_, ok := (<-m.updateChan)
		if !ok {
			break
		}
		// could do something to keep its data structure up to date
	}
}

// must be eventually closed with Close, or a goroutine leak will occure
func NewSpatialMap[T Spatialized]() *SpatialMap[T] {
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

	var handle *events.CancelHandle = nil
	if events := p.MovementEvents(); events != nil {
		channel := make(chan ObjectMovedEvent)
		handle = events.Subscribe(channel)

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
	}

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
			if o.subHandle != nil {
				o.subHandle.Cancel()
			}
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
		return t.GetAABB().Contains(point)
	})
}

func (m *SpatialMap[T]) GetAllIntersects(aabb AABB) []T {
	m.lock.RLock()
	defer m.lock.RUnlock()

	result := make([]T, 0)

	for _, obj := range m.objects {
		oaabb := obj.val.GetAABB()
		if oaabb.IsZero() || aabb.Intersects(oaabb) {
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
