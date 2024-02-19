package terrain

import (
	"fmt"
	"io"
	"sync"
	"time"

	. "creeps.heav.fr/geom"
	mathutils "creeps.heav.fr/math_utils"
	"github.com/rs/zerolog/log"
)

type Tilemap struct {
	// guards chunks
	chunkslock sync.RWMutex
	// generator can only be accessed (for read or write) with write lock on
	// chunkslock
	generator *ChunkGenerator
	chunks    map[Point]*TilemapChunk
}

func NewTilemap(generator *ChunkGenerator) Tilemap {
	return Tilemap{
		generator: generator,
		chunks:    make(map[Point]*TilemapChunk),
	}
}

// Use GetTile and SetTile instead, this is mainly for serialization and internal use
// if the chunk isn't generated this retuns nil
func (tilemap *Tilemap) GetChunk(chunkPos Point) *TilemapChunk {
	tilemap.chunkslock.RLock()
	defer tilemap.chunkslock.RUnlock()
	return tilemap.chunks[chunkPos]
}

// Retuns the already known chunk or generates it
func (tilemap *Tilemap) GenerateChunk(chunkPos Point) *TilemapChunk {
	// first try to read already existing chunk
	tilemap.chunkslock.RLock()
	if chunk := tilemap.chunks[chunkPos]; chunk != nil {
		tilemap.chunkslock.RUnlock()
		return chunk
	}
	tilemap.chunkslock.RUnlock()
	// no chunk = get write access
	tilemap.chunkslock.Lock()
	defer tilemap.chunkslock.Unlock()

	// there may be a race condition where the chunk was generated between lock
	// access so we must re-check for it
	if chunk := tilemap.chunks[chunkPos]; chunk != nil {
		return chunk
	}

	log.Trace().Any("pos", chunkPos).Msg("Generating chunk")
	start := time.Now()
	// Only now can we safely generate the chunk
	tilemap.chunks[chunkPos] = tilemap.generator.GenerateChunk(chunkPos)
	log.Debug().
		Any("pos", chunkPos).
		TimeDiff("took", time.Now(), start).
		Msg("Finished generating chunk")
	return tilemap.chunks[chunkPos]
}

// func (t Tilemap) GetOrCreateChunk(p Point) tilemapChunk {
// 	_, ok := t.chunks[p]
// 	if !ok {
// 		t.chunks[p] = tilemapChunk{}
// 	}
// 	return t.chunks[p]
// }

// Gets read access on the chunk and returns the value of the chunk
func (t *Tilemap) GetTile(p Point) Tile {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	return chunk.GetTile(Global2ChunkSubCoords(p))
}

// Gets write access on the chunk, the sets the given tile to the given value
// returning its previous value
func (t *Tilemap) SetTile(p Point, newVal Tile) Tile {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	return chunk.SetTile(Global2ChunkSubCoords(p), newVal)
}

func (t *Tilemap) ModifyTile(p Point, cb func (Tile) Tile) {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	chunk.ModifyTile(Global2ChunkSubCoords(p), cb)
}

func (t *Tilemap) PrintRegion(w io.Writer, from Point, upto Point) {
	min_x := mathutils.Min(from.X, upto.X)
	min_y := mathutils.Min(from.Y, upto.Y)
	max_x := mathutils.Max(from.X, upto.X)
	max_y := mathutils.Max(from.Y, upto.Y)

	for y := max_y-1; y > min_y; y-- {
		for x := min_x; x < max_x; x++ {
			t.GetTile(Point { X:x, Y:y }).Print(w)
		}
		fmt.Fprintln(w)
	}
}

// Returns a list of tiles in the given region
func (t *Tilemap) ObserveRegion(aabb AABB) []Tile {
	tiles := make([]Tile, 0)

	// FIXME: Lots of locks, but locking the correct chunks... too lazy rn
	for y := aabb.From.Y; y < aabb.Upto().Y; y++ {
		for x := aabb.From.X; x < aabb.Upto().X; x++ {
			tiles = append(tiles, t.GetTile(Point { X:x, Y:y }))
		}
	}

	return tiles
}
