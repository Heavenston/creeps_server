package terrain

import (
	"fmt"
	"io"
	"sync"
	"time"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	mathutils "github.com/heavenston/creeps_server/creeps_lib/math_utils"
	"github.com/rs/zerolog/log"
)

type Tilemap struct {
	// guards chunks
	chunkslock sync.RWMutex
	// generator can only be accessed (for read or write) with write lock on
	// chunkslock
	generator IGenerator
	chunks    map[Point]*Chunk
}

// generator can be nil in which case the default generator will be used
func NewTilemap(generator IGenerator) Tilemap {
	if generator == nil {
		generator = &DefaultGenerator{}
	}
	return Tilemap{
		generator: generator,
		chunks:    make(map[Point]*Chunk),
	}
}

// Use GetTile and SetTile instead, this is mainly for serialization and internal use
// if the chunk isn't generated this retuns nil
func (tilemap *Tilemap) GetChunk(chunkPos Point) *Chunk {
	tilemap.chunkslock.RLock()
	defer tilemap.chunkslock.RUnlock()
	return tilemap.chunks[chunkPos]
}

// Return the existing chunk at given pos or create a new not yet generated
// one
func (tilemap *Tilemap) CreateChunk(chunkPos Point) *Chunk {
	// first try to read already existing chunk
	tilemap.chunkslock.RLock()
	if chunk := tilemap.chunks[chunkPos]; chunk != nil {
		tilemap.chunkslock.RUnlock()
		return chunk
	}
	tilemap.chunkslock.RUnlock()

	tilemap.chunkslock.Lock()
	defer tilemap.chunkslock.Unlock()

	// there may be a race condition where the chunk was created between lock
	// access so we must re-check for it
	if chunk := tilemap.chunks[chunkPos]; chunk != nil {
		return chunk
	}

	chunk := tilemap.chunks[chunkPos]
	if chunk == nil {
		chunk = NewChunk(chunkPos)
		tilemap.chunks[chunkPos] = chunk
	}

	return chunk
}

// Like GetChunk but if it would return nil this will generate the chunk using
// the assigned generator.
func (tilemap *Tilemap) GenerateChunk(chunkPos Point) *Chunk {
	chunk := tilemap.CreateChunk(chunkPos)

	wc := chunk.WLock()

	if chunk.isGenerated.Load() {
		wc.UnLock()
		return chunk
	}

	log.Trace().Any("pos", chunkPos).Msg("Generating chunk")
	start := time.Now()
	tilemap.generator.GenerateChunk(&wc)
	log.Debug().
		Any("pos", chunkPos).
		TimeDiff("took", time.Now(), start).
		Msg("Finished generating chunk")

	// we do have write lock to chunk so we know no race-condition is possible
	chunk.isGenerated.Store(true)

	wc.UnLock()

	chunk.UpdatedEventProvider.Emit(GeneratedChunkEvent{})
	return chunk
}

// Generate the chunk if needed and calls terrain.Chunk.GetTile
func (t *Tilemap) GetTile(p Point) Tile {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	return chunk.GetTile(Global2ChunkSubCoords(p))
}

// Like GetTile but won't generate the chunk if unavailable
func (t *Tilemap) TryGetTile(p Point) Tile {
	chunk := t.CreateChunk(Global2ContainingChunkCoords(p))
	return chunk.GetTile(Global2ChunkSubCoords(p))
}

// Generate the chunk if needed and calls terrain.Chunk.SetTile
func (t *Tilemap) SetTile(p Point, newVal Tile) Tile {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	return chunk.SetTile(Global2ChunkSubCoords(p), newVal)
}

// Generate the chunk if needed and calls terrain.Chunk.ModifyTile
func (t *Tilemap) ModifyTile(p Point, cb func(Tile) Tile) {
	chunk := t.GenerateChunk(Global2ContainingChunkCoords(p))
	chunk.ModifyTile(Global2ChunkSubCoords(p), cb)
}

func (t *Tilemap) PrintRegion(w io.Writer, from Point, upto Point) {
	min_x := mathutils.Min(from.X, upto.X)
	min_y := mathutils.Min(from.Y, upto.Y)
	max_x := mathutils.Max(from.X, upto.X)
	max_y := mathutils.Max(from.Y, upto.Y)

	for y := max_y - 1; y > min_y; y-- {
		for x := min_x; x < max_x; x++ {
			t.GetTile(Point{X: x, Y: y}).Print(w)
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
			tiles = append(tiles, t.GetTile(Point{X: x, Y: y}))
		}
	}

	return tiles
}
