serve: build
	webappserve --root ./cmd/webapp/public

build: tidy
	go run ./cmd/pwasm/main.go ./cmd/webapp

test: fmt
	go test -race ./pkg/...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./pkg/...

tidy:
	go mod tidy

fmt:
	goimports -l -w ./pkg
