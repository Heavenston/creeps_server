
build_viewer_front:
	cd viewer/front; \
	npm install --include=dev; \
	npm run build;

dev: build_viewer_front
	go run .

.PHONY: dev build_viewer_front
