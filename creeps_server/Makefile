build: creeps_server

SRC_FILES != find . -name "*.go"

creeps_server: ${SRC_FILES}
	go build

clean:
	go clean

dev:
	go run -race . ${CREEPS_ARGS}

serve: creeps_server
	./creeps_server ${CREEPS_ARGS}

.PHONY: dev serve clean build
