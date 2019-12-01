export GO111MODULE=on
export GOBIN=${PWD}/bin

REVISION = $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags="-s -w -X 'github.com/blp1526/blkinfo.revision=$(REVISION)'"

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -rf bin/
	@echo

.PHONY: mod
mod:
	go mod tidy
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get github.com/goreleaser/goreleaser
	@echo

.PHONY: lint
lint: mod
	./bin/golangci-lint run ./...
	@echo

.PHONY: test
test: lint
	go test ./... -v --cover -race -covermode=atomic -coverprofile=coverage.txt
	@echo

.PHONY: build
build: test
	./bin/goreleaser release --rm-dist --snapshot
	@echo
