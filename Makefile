.PHONY: all
all: test

.PHONY: lint
lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	@echo
	golangci-lint run ./...
	@echo

.PHONY: test
test: lint
	go test ./... -v --cover -race -covermode=atomic -coverprofile=coverage.txt
	@echo

.PHONY: build
build: test
	go build -o bin/blkinfo ./cmd/blkinfo
	@echo
