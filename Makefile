
build_front:
	@make -C front build

dev:
	./parallel.sh dev

clean:
	make -C front clean

.PHONY: dev build_front clean
