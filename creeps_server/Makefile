creeps_server:
	go build

clean:
	go clean

dev:
	go run -race . ${CREEPS_ARGS}

serve: creeps_server
	./creeps_server ${CREEPS_ARGS}

.PHONY: dev serve clean
