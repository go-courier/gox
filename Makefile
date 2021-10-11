BASE_PATH=/

dev:
	BASE_PATH=$(BASE_PATH)  pnpx vite --host --config=vite.config.ts

build: pnpm.i
	BASE_PATH=$(BASE_PATH) pnpx vite --config=vite.config.ts build

build.gh-pages:
	$(MAKE) build -e BASE_PATH=/gox/

preview: build
	pnpx vite preview --host --port 8080 ./cmd/webapp

pnpm.i:
	pnpm i

test: tidy
	go test -race ./pkg/...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./pkg/...

tidy: fmt
	go mod tidy

fmt:
	goimports -l -w ./pkg
	