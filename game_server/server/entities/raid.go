package entities

import (
	"creeps.heav.fr/events"
	. "creeps.heav.fr/geom"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
	"creeps.heav.fr/spatialmap"
	"creeps.heav.fr/uid"
	"github.com/rs/zerolog/log"
)

type Raid struct {
	OwnerEntity
	ownerPlayerId uid.Uid

	server *Server
	id     uid.Uid

	campPosition Point
	lastRaiderSpawn int
}

func NewRaid(
	server *Server,
	ownerPlayerId uid.Uid,
) *Raid {
	raid := new(Raid)

	raid.server = server
	raid.id = uid.GenUid()

	raid.ownerPlayerId = ownerPlayerId

	return raid
}

func (raid *Raid) GetServer() *Server {
	return raid.server
}

func (raid *Raid) GetId() uid.Uid {
	return raid.id
}

// for IEntity
func (raid *Raid) GetAABB() AABB {
	return AABB {
		From: raid.campPosition,
		Size: Point { X: 1, Y: 1 },
	}
}

// for IEntity
func (raid *Raid) GetOwner() uid.Uid {
	return raid.ownerPlayerId
}

// for IEntity
func (raid *Raid) MovementEvents() *events.EventProvider[spatialmap.ObjectMovedEvent] {
	return nil
}

func (raid *Raid) Register() {
	var player, ok = raid.server.GetEntity(raid.ownerPlayerId).(*Player)
	if !ok || player == nil {
		log.Warn().Any("raid_id", raid.id).Any("owner_player", raid.ownerPlayerId).
			Msg("Invilid raid owner")
		return
	}

	raid.server.RegisterEntity(raid)

	raid.campPosition = raid.server.FindSpawnPoint(player.GetSpawnPoint(), 1, func(p Point) bool {
		found := false
		raid.server.ForEachEntity(func(entity IEntity) (shouldStop bool) {
			eplayer, ok := entity.(*Player)
			if !ok {
				return
			}
			found = eplayer.spawnPoint.Dist(p) < 15
			shouldStop = found
			return
		})
		return !found
	})
	raid.server.Tilemap().SetTile(raid.campPosition, terrain.Tile {
		Kind: terrain.TileRaiderCamp,
		Value: 0,
	})

	log.Info().Any("raid_id", raid.id).
		Any("point", raid.campPosition).
		Any("owner_player", raid.ownerPlayerId).
		Msg("Started raid")
}

func (raid *Raid) Unregister() {
	raid.ForEachEntities(func(entity IEntity) (shouldStop bool) {
		entity.Unregister()
		return
	})

	raid.server.Tilemap().SetTile(raid.campPosition, terrain.Tile {
		Kind: terrain.TileGrass,
		Value: 0,
	})

	raid.server.RemoveEntity(raid.id)
}

func (raid *Raid) Tick() {
	var player, ok = raid.server.GetEntity(raid.ownerPlayerId).(*Player)
	if !ok || player == nil {
		// player is dead ?
		raid.Unregister()
		return
	}

	currentTick := raid.server.Ticker().GetTickNumber()
	rate := raid.server.GetSetup().EnemyTickRate

	if currentTick - raid.lastRaiderSpawn < rate {
		return
	}

	raid.lastRaiderSpawn = currentTick

	raider := NewTurretUnit(raid.server, raid.id)
	raider.SetPosition(raid.campPosition)
	raider.Register()
}
