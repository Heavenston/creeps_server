package gameplay

import (
	. "lib.creeps.heav.fr/geom"
	"creeps.heav.fr/server"
	"creeps.heav.fr/server/entities"
	"lib.creeps.heav.fr/terrain"
)

// Spawns the given player town hall and everything it needs
func InitPlayer(
	srv *server.Server,
	player *entities.Player,
) (townhall Point, household Point, c1, c2 *entities.CitizenUnit) {
	player.Register()

	townhall = player.GetSpawnPoint()
	household = player.GetSpawnPoint().Plus(0, 1)

	srv.Tilemap().SetTile(townhall, terrain.Tile{
		Kind:  terrain.TileTownHall,
		Value: 0,
	})
	srv.Tilemap().SetTile(household, terrain.Tile{
		Kind:  terrain.TileHousehold,
		Value: 0,
	})

	player.AddTownHall(townhall)

	c1 = entities.NewCitizenUnit(srv, player.GetId())
	c1.SetPosition(household)
	c1.Register()
	c2 = entities.NewCitizenUnit(srv, player.GetId())
	c2.SetPosition(household)
	c2.Register()

	return
}
