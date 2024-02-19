package units

import (
	"fmt"

	"creeps.heav.fr/epita_api/model"
	. "creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
	. "creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
	"github.com/rs/zerolog/log"
)

// used by ApplyAction
func observe(unit IUnit, into *model.ObserveReport) {
	server := unit.GetServer()
	dist := unit.ObserveDistance() / 2
	// remainder upto is excluded
	aabb := AABB {
		From: unit.GetPosition().Minus(dist, dist),
		Size: Point {
			X: unit.ObserveDistance(),
			Y: unit.ObserveDistance(),
		},
	}

	units := server.Units().GetAllIntersects(aabb)
	into.Units = make([]model.Unit, 0, len(units))
	for _, ounit := range units {
		playerUsername := "server"
		if player := server.GetPlayerFromId(ounit.GetOwner()); player != nil {
			playerUsername = player.GetUsername()
		}

		into.Units = append(into.Units, model.Unit{
			OpCode:   ounit.GetOpCode(),
			Player:   playerUsername,
			Position: ounit.GetPosition(),
		})
	}

	tiles := server.Tilemap().ObserveRegion(aabb)
	into.Tiles = make([]uint16, 0, len(units))
	for _, tile := range tiles {
		into.Tiles = append(into.Tiles, uint16(tile.Kind)<<10|uint16(tile.Value))
	}
}

// used by ApplyAction
func refine(
	unit IUnit,
	requiredTile terrain.TileKind,
	from model.ResourceKind,
	to model.ResourceKind,
) (report model.IReport) {
	server := unit.GetServer()
	ownerId := unit.GetOwner()
	player := server.GetPlayerFromId(ownerId)

	if player == nil {
		log.Fatal().Msg("no player no refine")
	}

	tile := unit.GetServer().Tilemap().GetTile(unit.GetPosition())
	if tile.Kind != requiredTile {
		report = &model.ErrorReport{
			ErrorCode: "not-on-suitable-refinery",
			Error: fmt.Sprintf(
				"Wrong tile, found %v but require %v",
				tile.Kind, requiredTile,
			),
		}
		return
	}

	player.ModifyResources(func(resources model.Resources) model.Resources {
		ptr := resources.OfKind(from)
		if *ptr <= 0 {
			report = &model.ErrorReport{
				ErrorCode: "insufficient-funds",
				Error:     fmt.Sprintf("Missing one %s to make one %s", from, to),
			}
			return resources
		}

		(*ptr)--

		(*resources.OfKind(to))++

		report = &model.RefineReport{}
		return resources
	})

	return
}

// used by ApplyAction
func build(
	unit IUnit,
	opcode string,
	cost *model.CostResponse,
	target terrain.TileKind,
) (report model.IReport) {
	server := unit.GetServer()
	ownerId := unit.GetOwner()
	player := server.GetPlayerFromId(ownerId)
	tilemap := server.Tilemap()
	position := unit.GetPosition()

	if player == nil {
		log.Fatal().Msg("no player no build")
	}

	tilemap.ModifyTile(position, func(t terrain.Tile) terrain.Tile {
		if t.Kind != terrain.TileGrass {
			report = &model.ErrorReport{
				Error:     "Tile must be empty for construction",
				ErrorCode: "tile-occupied",
			}
			return t
		}

		could := false
		player.ModifyResources(func(resources model.Resources) model.Resources {
			if resources.EnoughFor(cost.Resources) < 1 {
				report = &model.ErrorReport{
					Error:     "Not enough resources",
					ErrorCode: "insufficient-funds",
				}
				return resources
			}

			could = true
			resources.Remove(cost.Resources)
			return resources
		})

		if !could {
			return t
		}

		t.Kind = target
		t.Value = 0

		if opcode == "household" {
			c1 := NewCitizenUnit(server, player.GetId())
			c1.SetPosition(position)
			server.RegisterUnit(c1)
			c2 := NewCitizenUnit(server, player.GetId())
			c2.SetPosition(position)
			server.RegisterUnit(c2)

			report = &model.BuildHouseHoldReport{
				BuildReport: model.BuildReport{
					Building: model.Building{
						OpCode:   opcode,
						Player:   69,
						Position: position,
					},
				},
				SpawnedCitizen1Id: c1.id,
				SpawnedCitizen2Id: c2.id,
			}
		} else {
			report = &model.BuildReport{
				Building: model.Building{
					OpCode:   opcode,
					Player:   69,
					Position: position,
				},
			}
		}

		return t
	})

	return
}

// used by ApplyAction
func spawn[T IUnit](
	unit IUnit,
	cost *model.CostResponse,
	precreatedUnit T,
) (report model.IReport) {
	server := unit.GetServer()
	ownerId := unit.GetOwner()
	player := server.GetPlayerFromId(ownerId)

	if player == nil {
		log.Fatal().Msg("no player no spawn")
	}

	could := false
	player.ModifyResources(func(inv model.Resources) model.Resources {
		if inv.EnoughFor(cost.Resources) < 1 {
			report = &model.ErrorReport{
				Error:     "Not enough resources",
				ErrorCode: "insufficient-funds",
			}
			return inv
		}

		could = true
		inv.Remove(cost.Resources)
		return inv
	})

	if !could {
		return
	}

	precreatedUnit.SetPosition(unit.GetPosition())
	server.RegisterUnit(precreatedUnit)

	playerUsername := "server"
	if player := unit.GetServer().GetPlayerFromId(precreatedUnit.GetId()); player != nil {
		playerUsername = player.GetUsername()
	}

	report = &model.SpawnReport{
		SpawnedUnitId: precreatedUnit.GetId(),
		SpawnedUnit: model.Unit{
			OpCode:   precreatedUnit.GetOpCode(),
			Player:   playerUsername,
			Position: precreatedUnit.GetPosition(),
		},
	}

	return
}

// called by unit in units/unit.go when the action is finished
func ApplyAction(action *Action, unit IUnit) (report model.IReport) {
	server := unit.GetServer()
	player := server.GetPlayerFromId(unit.GetOwner())

	oldPosition := unit.GetPosition()

	switch action.OpCode {
	case OpCodeMoveLeft:
		fallthrough
	case OpCodeMoveRight:
		fallthrough
	case OpCodeMoveUp:
		fallthrough
	case OpCodeMoveDown:
		var newPos Point
		unit.ModifyPosition(func(pos Point) Point {
			newPos = pos.Add(action.OpCode.MoveDirection())
			return newPos
		})
		mv := &model.MoveReport{
			NewPosition: newPos,
		}
		observe(unit, &mv.ObserveReport)
		report = mv
	case OpCodeObserve:
		mv := &model.ObserveReport{}
		observe(unit, mv)
		report = mv
	case OpCodeGather:
		maxInventorySize := server.GetSetup().MaxLoad
		server.Tilemap().ModifyTile(unit.GetPosition(), func(tile terrain.Tile) terrain.Tile {
			resKind := tile.Kind.GetResourceName()
			if resKind == "" {
				log.Trace().
					Str("player_username", player.GetUsername()).
					Any("player_id", player.GetId()).
					Any("tile", tile).
					Msg("Gather on non-resources")
				report = &model.ErrorReport{
					ErrorCode: "not-resource-tile",
				}
				return tile
			}

			unit.ModifyInventory(func(res model.Resources) model.Resources {
				size := res.Size()
				if size >= maxInventorySize {
					return res
				}

				took := mathutils.Min(
					maxInventorySize-res.Size(),
					int(tile.Value),
				)
				*res.OfKind(resKind) += took

				tile.Value -= uint8(took)
				if tile.Value == 0 {
					tile.Kind = terrain.TileGrass
				}

				report = &model.GatherReport{
					Resource:      resKind,
					Gathered:      took,
					ResourcesLeft: int(tile.Value),
				}
				return res
			})
			return tile
		})
	case OpCodeUnload:
		tile := server.Tilemap().GetTile(unit.GetPosition())
		if tile.Kind != terrain.TileTownHall {
			report = &model.ErrorReport{
				ErrorCode: "not-on-town-hall",
				Error:     "Cannot unload unless on a TownHall tile",
			}
			break
		}

		var credited model.Resources
		unit.ModifyInventory(func(inv model.Resources) model.Resources {
			credited = inv
			player.ModifyResources(func(res model.Resources) model.Resources {
				return res.Sum(inv)
			})
			return model.Resources{}
		})
		report = &model.UnloadReport{
			CreditedResources: credited,
		}
	case OpCodeFarm:
		// TODO: Verify if there is water around
		// (this is harder than you think unless tilemap has better lock guarantees than rn)
		server.Tilemap().ModifyTile(unit.GetPosition(), func(tile terrain.Tile) terrain.Tile {
			if tile.Kind != terrain.TileGrass {
				report = &model.ErrorReport{
					ErrorCode: "tile-occupied",
					Error:     "no-water-nearby",
				}
				return tile
			}

			tile = terrain.Tile{
				Kind: terrain.TileBush,
				// FIXME: HAHAHA magic value lol
				Value: 20,
			}
			report = &model.FarmReport{
				FoodQuantity: int(tile.Value),
			}
			return tile
		})
	case OpCodeDismantle:
		report = &model.DismantleReport{}
	case OpCodeUpgrade:
		if unit.IsUpgraded() {
			report = &model.ErrorReport{
				Error:     "Unit is already upgraded",
				ErrorCode: "unit-already-upgraded",
			}
			break
		}

		var could bool

		unit.ModifyInventory(func(res model.Resources) model.Resources {
			if res.EnoughFor(unit.GetUpgradeCosts().Resources) < 1 {
				could = false
				return res
			}
			res.Remove(unit.GetUpgradeCosts().Resources)
			could = true
			return res
		})

		if !could {
			report = &model.ErrorReport{
				Error:     "Not enough resources for upgrade",
				ErrorCode: "insufficient-funds",
			}
			break
		}

		report = &model.UpgradeReport{}
	case OpCodeRefineCopper:
		report = refine(unit, terrain.TileSmeltery, model.Rock, model.Copper)
	case OpCodeRefineWoodPlank:
		report = refine(unit, terrain.TileSawMill, model.Wood, model.WoodPlank)
	case OpCodeBuildTownHall:
		report = build(unit, "town-hall", &server.GetCosts().BuildTownHall, terrain.TileTownHall)
	case OpCodeBuildHousehold:
		report = build(unit, "household", &server.GetCosts().BuildHousehold, terrain.TileHousehold)
	case OpCodeBuildSawmill:
		report = build(unit, "sawmill", &server.GetCosts().BuildSawmill, terrain.TileSawMill)
	case OpCodeBuildSmeltery:
		report = build(unit, "smeltery", &server.GetCosts().BuildSmeltery, terrain.TileSmeltery)
	case OpCodeBuildRoad:
		report = build(unit, "road", &server.GetCosts().BuildRoad, terrain.TileRoad)
	case OpCodeSpawnTurret:
		report = spawn[*TurretUnit](
			unit,
			&server.GetCosts().SpawnTurret,
			NewTurretUnit(server, unit.GetOwner()),
		)
	case OpCodeSpawnBomberBot:
		report = spawn[*TurretUnit](
			unit,
			&server.GetCosts().SpawnTurret,
			NewTurretUnit(server, unit.GetOwner()),
		)
	case OpCodeFireTurret:
	case OpCodeFireBomberBot:
	}

	report.GetReport().ReportId = action.ReportId
	report.GetReport().OpCode = string(action.OpCode)
	report.GetReport().UnitId = unit.GetId()
	report.GetReport().UnitPosition = oldPosition
	report.GetReport().Status = "SUCCESS"
	report.GetReport().Login = player.GetUsername()
	if _, ok := report.(*model.ErrorReport); ok {
		report.GetReport().Status = "ERROR"
	}

	return
}
