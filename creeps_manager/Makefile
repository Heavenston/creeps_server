build: front_build creeps_manager

front_build:
	make -C front build

templ_files := $(shell find . -name "*.templ")
templ_out_files := $(templ_files:.templ=_templ.go)

${templ_out_files}: ${templ_files}
	templ generate

src_files := $(shell find . -name "*.go")

creeps_manager: ${src_files} ${templ_out_files}
	go build

clean:
	make -C front clean
	go clean
	rm -rf $(templ_out_files)

dev: front_build ${templ_out_files}
	go run -race . ${CREEPS_MANAGER_ARGS}

serve: front_build creeps_manager
	./creeps_manager ${CREEPS_MANAGER_ARGS}

.PHONY: dev serve clean build front_build templ_build
