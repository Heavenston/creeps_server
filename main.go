package main

import (
	"time"

	"creeps.heav.fr/geom"
	"creeps.heav.fr/server/terrain"
)

func main() {
	generator := terrain.NewChunkGenerator(time.Now().UnixMilli())
	tilemap := terrain.NewTilemap(generator)
	tilemap.GenerateChunk(geom.Point{X: 0, Y: 0})
	tilemap.GetChunk(geom.Point{}).Print()
}
