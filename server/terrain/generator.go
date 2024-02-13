package terrain

import (
	"math/rand"

	"creeps.heav.fr/geom"
	simplex "github.com/ojrac/opensimplex-go"
)

type patch struct {
	scale  float64
	thresh float64

	kind             TileKind
	defaultTileValue uint8

	noise simplex.Noise
}

type ChunkGenerator struct {
	rand *rand.Rand

	patchs []patch
}

func (gen *ChunkGenerator) newPath(scale float64, thresh float64, kind TileKind, defaultTileValue uint8) {
	var elm = patch {
		scale:  scale,
		thresh: thresh,

		kind:             kind,
		defaultTileValue: defaultTileValue,

		noise: simplex.New(gen.rand.Int63()),
	}
	gen.patchs = append(gen.patchs, elm)
}

func NewChunkGenerator(seed int64) *ChunkGenerator {
	g := new(ChunkGenerator)
	g.rand = rand.New(rand.NewSource(seed))
	g.patchs = make([]patch, 0, 0)

	g.newPath(1./6., 0.6, TileWater, 0)
	g.newPath(1./6., 0.6, TileStone, 10)
	g.newPath(1./1., 0.75, TileOil, 10)
	g.newPath(1./3., 0.5, TileTree, 10)
	g.newPath(1./3., 0.4, TileBush, 10)

	return g
}

func (gen *ChunkGenerator) sample(x int, y int) Tile {
	for _, patch := range gen.patchs {
		val := patch.noise.Eval2(float64(x)*patch.scale, float64(y)*patch.scale)
		if val > patch.thresh {
			return Tile{
				Kind:  patch.kind,
				Value: patch.defaultTileValue,
			}
		}
	}

	return Tile{
		Kind:  TileGrass,
		Value: 0,
	}
}

func (gen *ChunkGenerator) GenerateChunk(chunkPos geom.Point) *TilemapChunk {
	chunk := new(TilemapChunk)
	wcl := chunk.WLock()
	defer wcl.UnLock()

	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			point := geom.Point{X: x, Y: y}
			wcl.SetTile(point, gen.sample(x+chunkPos.X*ChunkSize, y+chunkPos.Y*ChunkSize))
		}
	}

	return chunk
}
