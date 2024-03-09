
manage:
	@export CREEPS_MANAGER_ARGS="-vv"; \
	make -C creeps_manager dev

clean:
	@./parallel.sh clean creeps_server creeps_manager

.PHONY: dev serve clean
