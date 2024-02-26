package terrain

import (
	"fmt"
	"io"
	"sync"

	"github.com/heavenston/creeps_server/creeps_lib/events"
	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	mathutils "github.com/heavenston/creeps_server/creeps_lib/math_utils"
)

type TilemapUpdateEvent struct {
	UpdatedPosition Point
	PreviousValue Tile
	NewValue Tile
}

type TilemapChunk struct {
	chunkPos Point
	// guards tiles
	tileslock sync.RWMutex
	tiles [ChunkTileCount]Tile
	UpdatedEventProvider events.EventProvider[TilemapUpdateEvent]
}

type readLockedChunk struct {
    unlocked bool
    chunk *TilemapChunk
}

type writeLockedChunk struct {
    unlocked bool
    chunk *TilemapChunk
}

const ChunkSize = 32
const ChunkTileCount = ChunkSize * ChunkSize

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

func NewChunk(pos Point) (chunk *TilemapChunk) {
	chunk = new(TilemapChunk)

	chunk.chunkPos = pos
	for i := range chunk.tiles {
		chunk.tiles[i] = Tile {
			Kind: TileUnknown,
			Value: 0,
		}
	}

	return
}

// the "chunk position" (world pos / chunkSize)
func (chunk *TilemapChunk) GetChunkPos() Point {
	return chunk.chunkPos
}

func (chunk *TilemapChunk) IsInBounds(subcoord Point) bool {
    return subcoord.X >= 0 && subcoord.X < ChunkSize ||
		subcoord.Y >= 0 || subcoord.Y < ChunkSize
}

func (chunk *TilemapChunk) tileIndex(subcoord Point) int {
    return subcoord.X + subcoord.Y * ChunkSize
}

// Waits for a read access lock on the chunk and returns the value of the tile
func (chunk *TilemapChunk) GetTile(subcoord Point) Tile {
	if !chunk.IsInBounds(subcoord) {
        panic("out of bound chunk tile access")
	}

	chunk.tileslock.RLock()
	defer chunk.tileslock.RUnlock()
	return chunk.tiles[chunk.tileIndex(subcoord)]
}

// Waits for a write access lock on the chunk and sets the given tile to the
// given value, and returns the previous value
func (chunk *TilemapChunk) SetTile(subcoord Point, newValue Tile) Tile {
	return chunk.ModifyTile(subcoord, func(t Tile) Tile {
		return newValue
	})
}

// Atomically modify the given tile
// returns the pervious value
func (chunk *TilemapChunk) ModifyTile(subcoord Point, cb func(Tile) Tile) Tile {
	if !chunk.IsInBounds(subcoord) {
        panic("out of bound chunk tile access")
	}

	chunk.tileslock.Lock()
	tileRef := &chunk.tiles[chunk.tileIndex(subcoord)]
    prevValue := *tileRef
	newValue := cb(prevValue)
	*tileRef = newValue
	chunk.tileslock.Unlock()

	if newValue != prevValue {
		chunk.UpdatedEventProvider.Emit(TilemapUpdateEvent{
			UpdatedPosition: subcoord,
			PreviousValue: prevValue,
			NewValue: newValue,
		})
	}

	return prevValue
}

func (chunk *TilemapChunk) Print(w io.Writer) {
    rlc := chunk.RLock()
	defer rlc.UnLock()
	for y := ChunkSize - 1; y > 0; y-- {
		for x := 0; x < ChunkSize; x++ {
			point := Point{X: x, Y: y}
			rlc.GetTile(point).Print(w)
		}
		fmt.Fprintln(w)
	}
}

func (chunk *TilemapChunk) RLock() readLockedChunk {
    chunk.tileslock.RLock()
    return readLockedChunk{
        chunk: chunk,
    }
}

func (rlc *readLockedChunk) UnLock() {
    rlc.chunk.tileslock.Unlock()
    rlc.unlocked = true
}

func (rlc *readLockedChunk) GetTile(subcoords Point) Tile {
    return rlc.chunk.tiles[rlc.chunk.tileIndex(subcoords)]
}

func (chunk *TilemapChunk) WLock() writeLockedChunk {
    chunk.tileslock.Lock()
    return writeLockedChunk{
        chunk: chunk,
    }
}

func (wlc *writeLockedChunk) UnLock() {
    wlc.chunk.tileslock.Unlock()
    wlc.unlocked = true
}

func (rlc *writeLockedChunk) GetTile(subcoords Point) Tile {
    return rlc.chunk.tiles[rlc.chunk.tileIndex(subcoords)]
}

func (rlc *writeLockedChunk) SetTile(subcoords Point, newVal Tile) {
    rlc.chunk.tiles[rlc.chunk.tileIndex(subcoords)] = newVal
}
