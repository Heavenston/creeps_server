package terrain

import "github.com/heavenston/creeps_server/creeps_lib/geom"

type IGenerator interface {
	GenerateChunk(chunkPos geom.Point) *Chunk
}

type DefaultGenerator struct {
	
}

func (gen *DefaultGenerator) GenerateChunk(chunkPos geom.Point) *Chunk {
	return NewChunk(chunkPos)
}
