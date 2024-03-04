package terrain

import (
	"io"

	"github.com/fatih/color"
	"github.com/heavenston/creeps_server/creeps_lib/model"
	"github.com/rs/zerolog/log"
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

func TileFromResource(res model.ResourceKind) TileKind {
	switch res {
	case model.Food:
		return TileBush
	case model.Oil:
		return TileOil
	case model.Rock:
		return TileStone
	case model.Wood:
		return TileTree
	case model.Copper:
		fallthrough
	case model.WoodPlank:
		return TileUnknown

	default:
		log.Warn().Msg("TileFromResource switch not complete")
	}
	return TileUnknown
}

func (kind TileKind) TileToResource() model.ResourceKind {
	switch kind {
	case TileBush:
		return model.Food
	case TileTree:
		return model.Wood
	case TileStone:
		return model.Rock
	case TileOil:
		return model.Oil
	}
	return ""
}

func (kind TileKind) CanRefine(resource model.ResourceKind) bool {
	switch resource {
	case model.Copper:
		return kind == TileSmeltery
	case model.WoodPlank:
		return kind == TileSawMill
	}
	return true
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

	case TileUnknown:
		c.Add(color.Reset)
		c.Fprint(w, "XX")
	}
}
