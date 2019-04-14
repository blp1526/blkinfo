export GO111MODULE=on
export GOBIN=${PWD}/bin

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -rf bin/
	@echo

.PHONY: lint
lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	@echo
	./bin/golangci-lint run ./...
	@echo

.PHONY: test
test: lint
	go test ./... -v --cover -race -covermode=atomic -coverprofile=coverage.txt
	@echo

.PHONY: build
build: test
	go build -o bin/blkinfo ./cmd/blkinfo
	@echo
