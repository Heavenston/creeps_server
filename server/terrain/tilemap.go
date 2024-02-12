package terrain

import (
	"fmt"

	. "creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
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

type TilemapChunk struct {
	tiles [ChunkTileCount]Tile
}

type Tilemap struct {
	generator *ChunkGenerator
	chunks    map[Point]*TilemapChunk
}

func (chunk *TilemapChunk) GetTile(subcoord Point) *Tile {
	if subcoord.X < 0 || subcoord.X >= ChunkSize ||
		subcoord.Y < 0 || subcoord.Y >= ChunkSize {
		return nil
	}

	return &chunk.tiles[subcoord.X+subcoord.Y*ChunkSize]
}

func (chunk *TilemapChunk) Print() {
	for y := ChunkSize - 1; y > 0; y-- {
		for x := 0; x < ChunkSize; x++ {
			point := Point{X: x, Y: y}
			chunk.GetTile(point).Print()
		}
		fmt.Println()
	}
}

func NewTilemap(generator *ChunkGenerator) Tilemap {
	return Tilemap{
		generator: generator,
		chunks:    make(map[Point]*TilemapChunk),
	}
}

// Gets the position of the chunk containing the given global position
// Ex with ChunkSize as 16:
//
// From {x = 5, y = 38} -> {x = 0, y = 3}
// From {x = -1, y = 0} -> {x = -1, y = 0}
func Global2ContainingChunkCoords(tile Point) Point {
	return Point{
		X: mathutils.FloorDivInt(tile.X, ChunkSize),
		Y: mathutils.FloorDivInt(tile.Y, ChunkSize),
	}
}

// Gets the position *in the chunk* of the given global tile position
// Ex with ChunkSize as 16:
//
// From {x = 5, y = 38} -> {x = 5, y = 5}
// From {x = -1, y = 0} -> {x = 15, y = 0}
func Global2ChunkSubCoords(tile Point) Point {
	return Point{
		X: mathutils.RemEuclidInt(tile.X, ChunkSize),
		Y: mathutils.RemEuclidInt(tile.Y, ChunkSize),
	}
}

func (tilemap *Tilemap) GetChunk(chunkPos Point) *TilemapChunk {
	return tilemap.chunks[chunkPos]
}

func (tilemap *Tilemap) GenerateChunk(chunkPos Point) *TilemapChunk {
	if chunk := tilemap.GetChunk(chunkPos); chunk != nil {
		return chunk
	}

	tilemap.chunks[chunkPos] = tilemap.generator.GenerateChunk(chunkPos)
	return tilemap.chunks[chunkPos]
}

// func (t Tilemap) GetOrCreateChunk(p Point) tilemapChunk {
// 	_, ok := t.chunks[p]
// 	if !ok {
// 		t.chunks[p] = tilemapChunk{}
// 	}
// 	return t.chunks[p]
// }

func (t *Tilemap) GetTile(p Point) Tile {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	return *chunk.GetTile(Global2ChunkSubCoords(p))
}

func (t *Tilemap) SetTile(p Point, newVal Tile) {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	*chunk.GetTile(Global2ChunkSubCoords(p)) = newVal
}

func (t *Tilemap) PrintRegion(from Point, upto Point) {
	min_x := mathutils.MinInt(from.X, upto.X)
	min_y := mathutils.MinInt(from.Y, upto.Y)
	max_x := mathutils.MaxInt(from.X, upto.X)
	max_y := mathutils.MaxInt(from.Y, upto.Y)

	for y := max_y-1; y > min_y; y-- {
		for x := min_x; x < max_x; x++ {
			t.GetTile(Point { X:x, Y:y }).Print()
		}
		fmt.Println()
	}
}
