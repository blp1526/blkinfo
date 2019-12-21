export GO111MODULE=on
export GOBIN=${PWD}/bin

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -rf bin/
	@echo

.PHONY: mod
mod:
	go mod tidy
	go get github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get github.com/goreleaser/goreleaser
	@echo

.PHONY: editorconfig
editorconfig: mod
	./bin/editorconfig-checker
	@echo

.PHONY: lint
lint: editorconfig
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
