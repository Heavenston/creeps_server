package terrain

import (
	"fmt"

	"github.com/fatih/color"
)

type TileKind uint8

const (
	TileGrass TileKind = iota
	TileWater
	TileStone
	TileBush
	TileTree
	TileOil

	TileTownHall
	TileHousehold
	TileRoad
	TileSawMill
	TileSmeltery
)

const ChunkSize = 32
const ChunkTileCount = ChunkSize * ChunkSize

type Tile struct {
	Kind  TileKind
	Value uint8
}

func (tile Tile) Print() {
	color.Set(color.BgGreen)
	color.Set(color.FgBlack)
	switch tile.Kind {
	case TileGrass:
		fmt.Print("  ")
	case TileWater:
		color.Set(color.BgHiBlue)
		color.Set(color.FgBlue)
		fmt.Print("~ ")
	case TileStone:
		color.Set(color.BgHiBlack)
		color.Set(color.FgBlack)
		fmt.Print("# ")
	case TileBush:
		color.Set(color.FgHiRed)
		fmt.Print(". ")
	case TileTree:
		color.Set(color.FgHiGreen)
		fmt.Print("T ")
	case TileOil:
		color.Set(color.FgBlack)
		fmt.Print("■ ")

	case TileTownHall:
		fmt.Print("TH")
	case TileHousehold:
		fmt.Print("HH")
	case TileRoad:
		fmt.Print("RO")
	case TileSawMill:
		fmt.Print("SM")
	case TileSmeltery:
		fmt.Print("SL")
	}
	color.Unset()
}
