package terrain

import (
	"io"

	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/fatih/color"
)

type TileKind uint8

// follows epita's values
// so only add new ones do not change existing ones
const (
	TileGrass TileKind = iota
	TileWater
	TileStone
	TileTree
	TileBush
	TileOil

	TileTownHall
	TileHousehold
	TileSmeltery
	TileSawMill
	TileRaiderCamp
	TileRaiderBorder
	TileRoad

	// Used when the tile is not generated
	TileUnknown
)

func (kind TileKind) GetResourceName() model.ResourceKind {
	switch kind {
	case TileBush:
		return model.Food
	case TileTree:
		return model.Wood
	case TileOil:
		return model.Oil
	}
	return ""
}

type Tile struct {
	Kind  TileKind
	Value uint8
}

func (tile Tile) Print(w io.Writer) {
	c := color.New()
	c.Add(color.BgGreen)
	c.Add(color.FgBlack)
	switch tile.Kind {
	case TileGrass:
		c.Fprint(w, "  ")
	case TileWater:
		c.Add(color.BgHiBlue)
		c.Add(color.FgBlue)
		c.Fprint(w, "~ ")
	case TileStone:
		c.Add(color.BgHiBlack)
		c.Add(color.FgBlack)
		c.Fprint(w, "# ")
	case TileBush:
		c.Add(color.FgHiRed)
		c.Fprint(w, ". ")
	case TileTree:
		c.Add(color.FgHiGreen)
		c.Fprint(w, "T ")
	case TileOil:
		c.Add(color.FgBlack)
		c.Fprint(w, "■ ")

	case TileTownHall:
		c.Fprint(w, "TH")
	case TileHousehold:
		c.Fprint(w, "HH")
	case TileRoad:
		c.Fprint(w, "RO")
	case TileSawMill:
		c.Fprint(w, "SM")
	case TileSmeltery:
		c.Fprint(w, "SL")
	}
}
