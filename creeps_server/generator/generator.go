package generator

import (
	"math/rand"

	. "github.com/heavenston/creeps_server/creeps_lib/geom"
	. "github.com/heavenston/creeps_server/creeps_lib/terrain"
	simplex "github.com/ojrac/opensimplex-go"
)

type patch struct {
	thresh float64

	kind             TileKind
	defaultTileValue uint8

	sample func(x float64, y float64) float64
}

type NoiseGenerator struct {
	rand *rand.Rand

	patchs []patch
}

// First argument is a scale that should be applied to the simplex noise
// second is the minimum value the noise must be for the tile to be applied
// (note: simplex noise goes from 0 to 1)
// so higher threshold means the tile is rarer, and the lower the scale the more
// 'blury' and big the patches will be
func (gen *NoiseGenerator) newPath(
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

func NewNoiseGenerator(seed int64) *NoiseGenerator {
	g := new(NoiseGenerator)
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

func (gen *NoiseGenerator) sample(x int, y int) Tile {
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

func (gen *NoiseGenerator) GenerateChunk(wc *WriteLockedChunk) {
	chunk := wc.GetChunk()
	pos := chunk.GetChunkPos()

	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			point := Point{X: x, Y: y}
			wc.SetTile(point, gen.sample(x+pos.X*ChunkSize, y+pos.Y*ChunkSize))
		}
	}
}

