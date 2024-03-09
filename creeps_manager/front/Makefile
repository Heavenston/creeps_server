build: dist

node_modules: package.json package-lock.json
	npm install --include=dev
	# make sure the folder is marked as up to date
	@touch -m node_modules

src_files := $(shell find src static -type f) $(shell find .. -name "*.templ")

dist: node_modules ${src_files} build.mjs tailwind.config.js
	npm run build

clean:
	rm -rf dist
	rm -rf .parcel-cache

.PHONY: build clean
