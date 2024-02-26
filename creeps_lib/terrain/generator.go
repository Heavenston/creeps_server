package terrain

import "github.com/heavenston/creeps_server/creeps_lib/geom"

type IGenerator interface {
	GenerateChunk(chunkPos geom.Point) *TilemapChunk
}

type DefaultGenerator struct {
	
}

func (gen *DefaultGenerator) GenerateChunk(chunkPos geom.Point) *TilemapChunk {
	return NewChunk(chunkPos)
}
