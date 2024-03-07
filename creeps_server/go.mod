module github.com/heavenston/creeps_server/creeps_server

go 1.22.0

replace github.com/heavenston/creeps_server/creeps_lib => ../creeps_lib

require (
	github.com/alecthomas/kong v0.8.1
	github.com/go-chi/chi/v5 v5.0.11
	github.com/gorilla/websocket v1.5.1
	github.com/heavenston/creeps_server/creeps_lib v0.1.0
	github.com/ojrac/opensimplex-go v1.0.2
	github.com/rs/zerolog v1.32.0
)

require (
	github.com/fatih/color v1.16.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
