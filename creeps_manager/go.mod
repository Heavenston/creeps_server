module github.com/Heavenston/creeps_server/creeps_manager

go 1.22.0

replace github.com/heavenston/creeps_server/creeps_lib => ../creeps_lib

require (
	github.com/a-h/templ v0.2.598
	github.com/ajg/form v1.5.1
	github.com/alecthomas/kong v0.8.1
	github.com/go-chi/chi/v5 v5.0.12
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/heavenston/creeps_server/creeps_lib v0.1.0
	github.com/rs/zerolog v1.32.0
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gorm v1.25.7
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
