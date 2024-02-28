package server

import (
	"sync"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/spatialmap"
	"github.com/heavenston/creeps_server/creeps_lib/uid"
	"github.com/rs/zerolog/log"
)

// entities do not necessarily have a position
// they can be units, players or other abstract entities that needs to be
// ticked
type IEntity interface {
	GetServer() *Server
	GetId() uid.Uid
	GetAABB() AABB
	// can return nil if the entity cannot move
	MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent]
	// if there is no real owner use uid.ServerId
	GetOwner() uid.Uid

	IsRegistered() bool

	Unregister()
	Register()
	
	// Ran each tick after being registered by the server
	Tick()
}

// embed the OwnerEntity helper struct instead of implementing it again
type IOwnerEntity interface {
	IEntity
	CopyEntityList() map[uid.Uid]IEntity
	ForEachEntities(func (entity IEntity) (shouldStop bool))
	HasEntity(id uid.Uid) bool
	AddEntity(entity IEntity)
	RemoveEntity(id uid.Uid) IEntity
	OwnedEntityCount() int
}

// helper struct that implements methods needed for IOwnerEntity
// intended to be embedded in structs that needs it
type OwnerEntity struct {
	entitesLock   sync.RWMutex
	ownedEntities map[uid.Uid]IEntity
}

func (e *OwnerEntity) InitOwnedEntities() {
	e.ownedEntities = make(map[uid.Uid]IEntity)
}

func (e *OwnerEntity) CopyEntityList() map[uid.Uid]IEntity {
	e.entitesLock.RLock()
	defer e.entitesLock.RUnlock()

	copy := make(map[uid.Uid]IEntity, len(e.ownedEntities))
	for k, v := range e.ownedEntities {
		copy[k] = v
	}
	return copy
}

func (e *OwnerEntity) ForEachEntities(cb func (entity IEntity) (shouldStop bool)) {
	e.entitesLock.RLock()
	defer e.entitesLock.RUnlock()

	for _, entity := range e.ownedEntities {
		if cb(entity) {
			break
		}
	}
}

func (e *OwnerEntity) HasEntity(id uid.Uid) bool {
	e.entitesLock.RLock()
	defer e.entitesLock.RUnlock()

	return e.ownedEntities[id] != nil
}

func (e *OwnerEntity) AddEntity(entity IEntity) {
	e.entitesLock.Lock()
	defer e.entitesLock.Unlock()

	if e.ownedEntities[entity.GetId()] != nil {
		log.Warn().
			Str("id", string(entity.GetId())).
			Msg("attempted to add entity with the same id as an already owned one")
		return
	}

	e.ownedEntities[entity.GetId()] = entity
}

func (e *OwnerEntity) RemoveEntity(id uid.Uid) IEntity {
	e.entitesLock.Lock()
	defer e.entitesLock.Unlock()

	entity := e.ownedEntities[id]
	delete(e.ownedEntities, id)
	return entity
}

func (f *OwnerEntity) OwnedEntityCount() int {
	return len(f.ownedEntities)
}
