test: tidy
	go test -race ./...

cover: tidy
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

tidy: fmt
	go mod tidy

fmt:
	goimports -l -w ./pkg

dep:
	go get -u ./...