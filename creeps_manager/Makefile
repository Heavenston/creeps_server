build: front_build creeps_manager

SRC_FILES != find . -name "*.go"

creeps_manager: ${SRC_FILES}
	go build

clean:
	make -C front clean
	go clean

dev: front_build
	go run -race . ${CREEPS_MANAGER_ARGS}

serve: front_build creeps_manager
	./creeps_manager ${CREEPS_MANAGER_ARGS}

.PHONY: dev serve clean build front_build
