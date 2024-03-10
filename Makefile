
manage:
	if [ -f .env ]; then export $$(cat .env | xargs); fi; \
	make -C creeps_manager dev

build:
	make -C creeps_manager build
	make -C creeps_server build

clean:
	make -C creeps_manager clean
	make -C creeps_server clean

.PHONY: dev serve clean
