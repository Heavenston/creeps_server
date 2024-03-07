
manage:
	@export CREEPS_MANAGER_ARGS="-vv"; \
	./parallel.sh dev creeps_manager front

dev:
	@export CREEPS_ARGS="-v"; \
	./parallel.sh dev creeps_server front

trace:
	@export CREEPS_ARGS="-vv"; \
	./parallel.sh dev creeps_server front

serve:
	@./parallel.sh serve creeps_server front

clean:
	@./parallel.sh clean creeps_server creeps_manager front

.PHONY: dev serve clean
