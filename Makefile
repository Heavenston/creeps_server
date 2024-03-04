
dev:
	@export CREEPS_ARGS="-l debug"; \
	./parallel.sh dev

trace:
	@export CREEPS_ARGS="-l trace"; \
	./parallel.sh dev

serve:
	@./parallel.sh serve

clean:
	@./parallel.sh clean

.PHONY: dev serve clean
