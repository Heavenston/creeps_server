package creepsclientlib

import (
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/heavenston/creeps_server/creeps_lib/terrain"
)

func registerTiles(tp *terrain.Tilemap, tiles []uint16) {

}

func registerBuilding(tp *terrain.Tilemap, buildReport *model.BuildReport) {
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

	tp.SetTile(buildReport.UnitPosition, terrain.Tile{
		Kind:  tileKind,
		Value: uint8(buildReport.Building.Player),
	})
}

// Applies the tile modifications implied by the given report
func RegisterReport(tp *terrain.Tilemap, report model.IReport) {
	switch casted := report.(type) {
	case *model.ObserveReport:
		registerTiles(tp, casted.Tiles)
	case *model.MoveReport:
		registerTiles(tp, casted.Tiles)
	case *model.GatherReport:
		if casted.ResourcesLeft == 0 {
			tp.SetTile(casted.UnitPosition, terrain.Tile{
				Kind:  terrain.TileGrass,
				Value: 0,
			})
		} else {
			tp.SetTile(casted.UnitPosition, terrain.Tile{
				Kind:  terrain.TileFromResource(casted.Resource),
				Value: uint8(casted.ResourcesLeft),
			})
		}
	case *model.FarmReport:
		tp.SetTile(casted.UnitPosition, terrain.Tile{
			Kind:  terrain.TileFromResource(model.Food),
			Value: uint8(casted.FoodQuantity),
		})
	case *model.BuildReport:
		registerBuilding(tp, casted)
	case *model.BuildHouseHoldReport:
		registerBuilding(tp, &casted.BuildReport)
	}
}
