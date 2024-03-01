module github.com/heavenston/creeps_server/creeps_client_lib

go 1.22.0

replace github.com/heavenston/creeps_server/creeps_lib => ../creeps_lib

require github.com/heavenston/creeps_server/creeps_lib v0.0.0

require (
	github.com/fatih/color v1.16.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
