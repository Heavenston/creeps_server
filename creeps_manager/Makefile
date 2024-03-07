build: creeps_manager

SRC_FILES != find . -name "*.go"

creeps_manager: ${SRC_FILES}
	go build

clean:
	go clean

dev:
	go run -race . ${CREEPS_MANAGER_ARGS}

serve: creeps_manager
	./creeps_manager ${CREEPS_MANAGER_ARGS}

.PHONY: dev serve clean build
