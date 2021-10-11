dev:
	pnpx vite --host --config=vite.config.ts

BASE_PATH=/

build: pnpm.i
	BASE_PATH=$(BASE_PATH) pnpx vite --config=vite.config.ts build

build.gh-pages:
	$(MAKE) build -e BASE_PATH=/gox/

preview: build
	pnpx vite preview --host --port 8080 ./cmd/webapp

pnpm.i:
	pnpm i

test: fmt
	go test -race ./pkg/...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./pkg/...

tidy:
	go mod tidy

fmt:
	goimports -l -w ./pkg
