build: dist

node_modules: package.json package-lock.json
	npm install --include=dev
	# make sure the folder is marked as up to date
	@touch -m node_modules

src_files := $(shell find src -name "*.ts" -o -name "*.html" -o -name "*.css")

dist: node_modules ${src_files}
	npm run build > /dev/null

dev: serve

serve: dist
	npm run serve

clean:
	rm -rf dist
	rm -rf .parcel-cache

.PHONY: build clean serve
