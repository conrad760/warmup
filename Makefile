BINARY    := warmup
GOFILES   := $(wildcard *.go)

.PHONY: build install clean fmt lint vet test run stats help

## build: compile the binary (default)
build: $(BINARY)

$(BINARY): $(GOFILES) go.mod go.sum
	go build -o $(BINARY) .

## install: install to $GOPATH/bin
install:
	go install .

## run: build and launch a session
run: build
	./$(BINARY)

## stats: build and show lifetime stats
stats: build
	./$(BINARY) --stats

## fmt: format all Go source files
fmt:
	gofmt -w .

## lint: run staticcheck (provided by flake devShell)
lint: vet
	staticcheck ./...

## vet: run go vet
vet:
	go vet ./...

## test: run tests
test:
	go test ./...

## clean: remove compiled binaries
clean:
	rm -f $(BINARY)
	go clean

## help: show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
