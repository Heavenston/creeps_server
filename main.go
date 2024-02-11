package main

import (
	"time"

	"creeps.heav.fr/geom"
	"creeps.heav.fr/server/terrain"
)

func main() {
	generator := terrain.NewChunkGenerator(time.Now().UnixMilli())
	tilemap := terrain.NewTilemap(generator)
	tilemap.GenerateChunk(geom.Point{X: 0, Y: 2}).Print()
	tilemap.GenerateChunk(geom.Point{X: 0, Y: 1}).Print()
	tilemap.GenerateChunk(geom.Point{X: 0, Y: 0}).Print()
}
