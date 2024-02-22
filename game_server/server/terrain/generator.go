package terrain

import (
	"math/rand"

	"creeps.heav.fr/geom"
	simplex "github.com/ojrac/opensimplex-go"
)

type patch struct {
	thresh float64

	kind             TileKind
	defaultTileValue uint8

	sample func(x float64, y float64) float64
}

type ChunkGenerator struct {
	rand *rand.Rand

	patchs []patch
}

// First argument is a scale that should be applied to the simplex noise
// second is the minimum value the noise must be for the tile to be applied
// (note: simplex noise goes from 0 to 1)
// so higher threshold means the tile is rarer, and the lower the scale the more
// 'blury' and big the patches will be
func (gen *ChunkGenerator) newPath(
	scale float64,
	thresh float64,
	kind TileKind,
	defaultTileValue uint8,
) {
	noise := simplex.NewNormalized(gen.rand.Int63())
	
	var elm = patch {
		thresh: thresh,

		kind:             kind,
		defaultTileValue: defaultTileValue,

		sample: func(x, y float64) float64 {
			v1 := noise.Eval2(x*scale, y*scale)
			return v1
		},
	}
	gen.patchs = append(gen.patchs, elm)
}

func NewChunkGenerator(seed int64) *ChunkGenerator {
	g := new(ChunkGenerator)
	g.rand = rand.New(rand.NewSource(seed))
	g.patchs = make([]patch, 0, 0)

	// see newPath docs
	g.newPath(1./6., 0.85, TileWater, 0)
	g.newPath(1./5., 0.85, TileStone, 10)
	g.newPath(1./1., 0.95, TileOil, 10)
	g.newPath(1./3., 0.80, TileTree, 10)
	g.newPath(1./3., 0.75, TileBush, 10)

	return g
}

func (gen *ChunkGenerator) sample(x int, y int) Tile {
	for _, patch := range gen.patchs {
		val := patch.sample(float64(x), float64(y))
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

	chunk.chunkPos = chunkPos

	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			point := geom.Point{X: x, Y: y}
			wcl.SetTile(point, gen.sample(x+chunkPos.X*ChunkSize, y+chunkPos.Y*ChunkSize))
		}
	}

	return chunk
}
