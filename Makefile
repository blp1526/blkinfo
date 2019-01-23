.PHONY: all
all: test

.PHONY: gometalinter
gometalinter:
	go get github.com/alecthomas/gometalinter
	gometalinter --install
	@echo
	gometalinter --config .gometalinter.json ./...
	@echo

.PHONY: test
test: gometalinter
	go test ./... -v --cover -race -covermode=atomic -coverprofile=coverage.txt
	@echo
