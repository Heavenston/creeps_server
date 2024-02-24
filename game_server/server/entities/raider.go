package entities

import (
	"sync"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
	"creeps.heav.fr/uid"
	"github.com/rs/zerolog/log"
)

type RaiderUnit struct {
	unit
	lock   sync.RWMutex

	owner uid.Uid

	target Point
}

func NewRaiderUnit(server *Server, owner uid.Uid, target Point) *RaiderUnit {
	raider := new(RaiderUnit)
	raider.unitInit(server)
	raider.this = raider

	raider.target = target
	raider.owner = owner

	return raider
}

// for the extendedUnit interface
func (raider *RaiderUnit) getUnit() *unit {
	return &raider.unit
}

func (raider *RaiderUnit) GetOpCode() string {
	return "raider"
}

func (raider *RaiderUnit) GetUpgradeCosts() *model.CostResponse {
	return nil
}

func (raider *RaiderUnit) GetOwner() uid.Uid {
	return raider.owner
}

func (raider *RaiderUnit) GetTarget() Point {
	return raider.target
}

func (raider *RaiderUnit) StartAction(action *Action) error {
	err := raider.startAction(action, []ActionOpCode {
		OpCodeMoveDown,
		OpCodeMoveUp,
		OpCodeMoveLeft,
		OpCodeMoveRight,
	})
	if err != nil {
		return err
	}
	return nil
}

func (raider *RaiderUnit) Tick() {
	raider.lock.Lock()
	defer raider.lock.Unlock()

	raider.tick()

	owner := raider.server.GetEntityOwner(raider.id)
	if owner == nil {
		log.Warn().
			Str("raider_id", string(raider.id)).
			Str("owner_id", string(raider.owner)).
			Any("target", raider.target).
			Msg("RAIDER: Could not find my owner so imma kms")
		raider.Unregister()
		return
	}

	position := raider.GetPosition()

	foundAndDestroy := false
	raider.server.Tilemap().ModifyTile(position, func(t terrain.Tile) terrain.Tile {
		destroy := t.Kind == terrain.TileRoad ||
				   t.Kind == terrain.TileHousehold ||
				   t.Kind == terrain.TileSawMill ||
				   t.Kind == terrain.TileTownHall ||
				   t.Kind == terrain.TileSmeltery
		foundAndDestroy = foundAndDestroy || destroy
		if destroy {
			t.Kind = terrain.TileGrass
			t.Value = 0
		}
		return t
	})

	for _, entity := range raider.server.Entities().GetAllIntersects(raider.GetAABB()) {
		_, isC := entity.(*CitizenUnit)
		_, isT := entity.(*TurretUnit)
		_, isB := entity.(*BomberBotUnit)
		if isC || isT || isB {
			foundAndDestroy = true
			entity.Unregister()
		}
	}

	if foundAndDestroy {
		raider.SetDead()
		return
	}

	// busy = do nothing
	if action := raider.GetLastAction(); action != nil && !action.Finised.Load() {
		return
	}

	if raider.target == position {
		raider.SetDead()
		return
	}

	diff := raider.target.Sub(position)
	newAction := new(Action)

	if mathutils.AbsInt(diff.X) > mathutils.AbsInt(diff.Y) {
		if diff.X < 0 {
			newAction.OpCode = OpCodeMoveLeft
		} else {
			newAction.OpCode = OpCodeMoveRight
		}
	} else {
		if diff.Y < 0 {
			newAction.OpCode = OpCodeMoveDown
		} else {
			newAction.OpCode = OpCodeMoveUp
		}
	}

	err := raider.StartAction(newAction)
	if err != nil {
		log.Warn().
			Any("action", newAction).
			Err(err).
			Msg("[RAIDER] Could not start action")
	}
}
