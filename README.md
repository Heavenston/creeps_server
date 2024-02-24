# Creeps

Implementation in Go of the very famous Creeps game

## Development (/ the only way to use it for now lol)

You need go as well as nodejs/npm installed on your computer.
(TODO: minimum versions ?)

Run
```bash
make dev
```
to start the develpment environment.

A web server is opened at port `1234`, while the epita-compatible creeps
"rest" api is opened at port `1664`.

## Todo

- [ ] Game manager
	- [ ] Add a list of games to the main menu (just be able to add ips manually at least)
	- [ ] Be able to create/join games from a master server
- [x] Map generation (not at all like the real one)
- [x] Cli
	- [ ] Accept config files (using spf13/viper)
- [ ] Game viewer
	- [x] Load chunks and units from the server
	- [x] Render textures !
	- [ ] TEXTURE PACKKKKS
- [ ] Actions
	- [x] Movements
	- [x] Observe
	- [x] Gather (see known bugs)
	- [x] Unload
	- [x] Farm
	- [ ] Dismantle
	- [x] Upgrade
	- [x] Refine
	- [x] Build
	- [x] Spawn
	- [x] Fire (turret)
	- [x] Fire (bomber-bot)
- [x] Enemies
	- [ ] Make them get stronger and stronger (whatever that means)
- [ ] Garbage collector
- [ ] LOTS OF TESTING (and tests? lol)
- [ ] More techtree stuff like machine guns or nuclear bombs (really important) for pvp

### Known bugs

- Units seems to teleports (unitMoved packets are either not sent or missed)
  (maybe fixed)
