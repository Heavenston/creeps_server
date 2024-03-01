package terrain

import "github.com/heavenston/creeps_server/creeps_lib/geom"

type IGenerator interface {
	GenerateChunk(into *WriteLockedChunk)
}

type DefaultGenerator struct {
}

func (gen *DefaultGenerator) GenerateChunk(wc *WriteLockedChunk) {
	for y := 0; y < ChunkSize; y++ {
		for x := 0; x < ChunkSize; x++ {
			wc.SetTile(geom.Point{X: x, Y: y}, Tile{
				Kind: TileUnknown,
			})
		}
	}
}
