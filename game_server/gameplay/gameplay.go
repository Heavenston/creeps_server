package gameplay

import (
	"creeps.heav.fr/server"
	"creeps.heav.fr/server/terrain"
	"creeps.heav.fr/units"
	. "creeps.heav.fr/geom"
)

// Spawns the given player town hall and everything it needs
func InitPlayer(
	srv *server.Server,
	player *server.Player,
) (townhall Point, household Point, c1, c2 *units.CitizenUnit) {
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

	c1 = units.NewCitizenUnit(srv, player.GetId())
	c1.SetPosition(household)
	c1.Register()
	c2 = units.NewCitizenUnit(srv, player.GetId())
	c2.SetPosition(household)
	c2.Register()

	return
}
