
manage:
	@export CREEPS_MANAGER_ARGS="-vv"; \
	./parallel.sh dev creeps_manager

build:
	./parallel.sh build creeps_manager

clean:
	@./parallel.sh clean creeps_server creeps_manager

.PHONY: dev serve clean
