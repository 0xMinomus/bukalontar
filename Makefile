.PHONY: install dev dev-api web-build api-build build preview test lint clean

install:
	npm install

dev:
	npm run dev

dev-api:
	cd apps/api && go run ./cmd/server

web-build:
	npm run build

api-build: web-build
	rm -rf apps/api/internal/webui/dist
	mkdir -p apps/api/internal/webui/dist
	cp -r apps/web/dist/. apps/api/internal/webui/dist/
	touch apps/api/internal/webui/dist/.gitkeep
	mkdir -p bin
	cd apps/api && go build -o ../../bin/bukalontar ./cmd/server

build: api-build

preview:
	npm run preview

test:
	cd apps/api && go test ./...

lint:
	npm run lint
	cd apps/api && go vet ./...

clean:
	rm -rf apps/web/dist
	rm -rf bin
	find apps/api/internal/webui/dist -mindepth 1 ! -name .gitkeep -delete
