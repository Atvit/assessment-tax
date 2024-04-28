.PHONY: build
build:
	go build -o bin/server server/server.go

.PHONY: start
start:
	go run main.go

.PHONY: test
test:
	go test -v $$(go list ./... | grep -v /mocks/) -cover

.PHONY: coverage
coverage:
	go test -v $$(go list ./... | grep -v /mocks/) -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: vet
vet:
	go vet ./...

.PHONY: mockgen
mockery:
	mockery --all --dir=./internals/setting --output=./mocks/setting --case=underscore --outpkg=mocks