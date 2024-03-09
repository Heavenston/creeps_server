
manage:
	@export CREEPS_MANAGER_ARGS="-vv"; \
	./parallel.sh dev creeps_manager front

clean:
	@./parallel.sh clean creeps_server creeps_manager front

.PHONY: dev serve clean
