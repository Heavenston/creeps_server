package main

import (
	"creeps.heav.fr/geom"
	"creeps.heav.fr/server"
)

func main() {
	generator := server.NewChunkGenerator(5)
	tilemap := server.NewTilemap(generator)
	tilemap.GenerateChunk(geom.Point{X: 0, Y: 0})
	tilemap.GetChunk(geom.Point{}).Print()
}
