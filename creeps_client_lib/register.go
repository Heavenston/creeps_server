package creepsclientlib

import (
	"math"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
)

func registerTiles(client *Client, pos Point, tiles []uint16) {
	size := int(math.Ceil(math.Sqrt(float64(len(tiles)))))
	for i, val := range tiles {
		x := i % size
		y := i / size
		client.tilemap.Load().SetTile(pos.Plus(x-size/2, y-size/2), terrain.Tile{
			Kind:  terrain.TileKind(val >> 10),
			Value: uint8(val & 0x3F),
		})
	}
}

func registerBuilding(client *Client, buildReport *model.BuildReport) {
	var tileKind terrain.TileKind

	// FIXME: Make an 'enum' somewhere
	switch buildReport.Building.OpCode {
	case "town-hall":
		tileKind = terrain.TileTownHall
	case "household":
		tileKind = terrain.TileHousehold
	case "sawmill":
		tileKind = terrain.TileSawMill
	case "smeltery":
		tileKind = terrain.TileSmeltery
	case "road":
		tileKind = terrain.TileRoad
	default:
		return
	}

	client.tilemap.Load().SetTile(buildReport.UnitPosition, terrain.Tile{
		Kind:  tileKind,
		Value: uint8(buildReport.Building.Player),
	})
}

// Applies the tile and resource modifications implied by the given report
func RegisterReport(client *Client, report model.IReport) {
	switch casted := report.(type) {
	case *model.ObserveReport:
		registerTiles(client, casted.UnitPosition, casted.Tiles)
	case *model.MoveReport:
		registerTiles(client, casted.NewPosition, casted.Tiles)
	case *model.RefineReport:
		to := casted.OpCode.RefineEndResult()
		client.playerResources.Modify(func(r model.Resources) model.Resources {
			*r.OfKind(to)++
			return r
		})
	case *model.GatherReport:
		client.UnitResources(casted.UnitId).Modify(func(r model.Resources) model.Resources {
			*r.OfKind(casted.Resource) += casted.Gathered
			return r
		})

		if casted.ResourcesLeft == 0 {
			client.tilemap.Load().SetTile(casted.UnitPosition, terrain.Tile{
				Kind:  terrain.TileGrass,
				Value: 0,
			})
		} else {
			client.tilemap.Load().SetTile(casted.UnitPosition, terrain.Tile{
				Kind:  terrain.TileFromResource(casted.Resource),
				Value: uint8(casted.ResourcesLeft),
			})
		}
	case *model.FarmReport:
		client.UnitResources(casted.UnitId)

		client.tilemap.Load().SetTile(casted.UnitPosition, terrain.Tile{
			Kind:  terrain.TileFromResource(model.Food),
			Value: uint8(casted.FoodQuantity),
		})
	case *model.BuildReport:
		registerBuilding(client, casted)
	case *model.BuildHouseHoldReport:
		registerBuilding(client, &casted.BuildReport)
	}
}
