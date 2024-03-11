
manager_build: tygo
	make -C creeps_manager build

server_build:
	make -C creeps_server build

build: tygo manager_build server_build

# Yes it runs every time, no i don't feel bad about it
tygo:
	tygo generate

manage: server_build tygo
	if [ -f .env ]; then export $$(cat .env | xargs); fi; \
	make -C creeps_manager dev

clean:
	make -C creeps_manager clean
	make -C creeps_server clean
	${RM} /creeps_manager/front/src/models/epita.ts
	${RM} /creeps_manager/front/src/models/viewer.ts

.PHONY: dev serve clean tygo manager_build server_build
