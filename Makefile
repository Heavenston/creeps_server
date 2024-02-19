
build_front:
	@make -C front build

dev:
	./dev.sh

clean:
	make -C front clean

.PHONY: dev build_front clean
