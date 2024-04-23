.PHONY: build
build:
	go build -o bin/server server/server.go

.PHONY: start
start:
	go run main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: coverage
coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out