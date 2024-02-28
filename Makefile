
dev:
	@export CREEPS_ARGS="-v"; \
	./parallel.sh dev

trace:
	@export CREEPS_ARGS="-vv"; \
	./parallel.sh dev

serve:
	@./parallel.sh serve

clean:
	@./parallel.sh clean

.PHONY: dev serve clean
