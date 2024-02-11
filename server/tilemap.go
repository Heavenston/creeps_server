package server

import (
	"fmt"

	"creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
)

type TileKind int

const (
	TileGrass TileKind = iota
	TileWater
	TileStone
	TileBush
	TileTree
	TileOil
)

const ChunkSize = 16
const ChunkTileCount = ChunkSize * ChunkSize

type Tile struct {
	kind  TileKind
	value int
}

type TilemapChunk struct {
	tiles [ChunkTileCount]Tile
}

type Tilemap struct {
	generator *ChunkGenerator
	chunks    map[geom.Point]*TilemapChunk
}

func (chunk *TilemapChunk) GetTile(subcoord geom.Point) *Tile {
	if subcoord.X < 0 || subcoord.X >= ChunkSize ||
		subcoord.Y < 0 || subcoord.Y >= ChunkSize {
		return nil
	}

	return &chunk.tiles[subcoord.X+subcoord.Y*ChunkSize]
}

func (chunk *TilemapChunk) Print() {
	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			point := geom.Point{X: x, Y: y}
			switch kind := chunk.GetTile(point).kind; kind {
			case TileGrass:
				fmt.Print(" ")
			case TileWater:
				fmt.Print("~")
			case TileStone:
				fmt.Print("#")
			case TileBush:
				fmt.Print(".")
			case TileTree:
				fmt.Print("T")
			case TileOil:
				fmt.Print("O")
			}
		}
		fmt.Println()
	}
}

func NewTilemap(generator *ChunkGenerator) Tilemap {
	return Tilemap{
		generator: generator,
		chunks:    make(map[geom.Point]*TilemapChunk),
	}
}

// Gets the position of the chunk containing the given global position
// Ex with ChunkSize as 16:
//
// From {x = 5, y = 38} -> {x = 0, y = 3}
// From {x = -1, y = 0} -> {x = -1, y = 0}
func Global2ContainingChunkCoords(tile geom.Point) geom.Point {
	return geom.Point{
		X: mathutils.FloorDivInt(tile.X, ChunkSize),
		Y: mathutils.FloorDivInt(tile.Y, ChunkSize),
	}
}

// Gets the position *in the chunk* of the given global tile position
// Ex with ChunkSize as 16:
//
// From {x = 5, y = 38} -> {x = 5, y = 5}
// From {x = -1, y = 0} -> {x = 15, y = 0}
func Global2ChunkSubCoords(tile geom.Point) geom.Point {
	return geom.Point{
		X: mathutils.RemEuclidInt(tile.X, ChunkSize),
		Y: mathutils.RemEuclidInt(tile.Y, ChunkSize),
	}
}

func (tilemap Tilemap) GetChunk(chunkPos geom.Point) *TilemapChunk {
	return tilemap.chunks[chunkPos]
}

func (tilemap Tilemap) GenerateChunk(chunkPos geom.Point) {
	if tilemap.GetChunk(chunkPos) != nil {
		return
	}

	tilemap.chunks[chunkPos] = tilemap.generator.GenerateChunk(chunkPos)
}

// func (t Tilemap) GetOrCreateChunk(p geom.Point) tilemapChunk {
// 	_, ok := t.chunks[p]
// 	if !ok {
// 		t.chunks[p] = tilemapChunk{}
// 	}
// 	return t.chunks[p]
// }

func (t Tilemap) GetTile(p geom.Point) *Tile {
	chunk := t.GetChunk(Global2ContainingChunkCoords(p))
	if chunk == nil {
		return nil
	}
	return chunk.GetTile(Global2ChunkSubCoords(p))
}
