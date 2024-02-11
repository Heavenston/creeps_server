package server

import "creeps.heav.fr/geom"

type ChunkGenerator struct {
}

func NewChunkGenerator(seed uint64) *ChunkGenerator {
	g := new(ChunkGenerator)
	return g
}

func (gen *ChunkGenerator) GenerateChunk(chunkPos geom.Point) *TilemapChunk {
	chunk := new(TilemapChunk)

	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			point := geom.Point{X: x, Y: y}
			chunk.GetTile(point).kind = TileGrass
		}
	}

	return chunk
}
