MOD_NAME := uptimer-cli
BIN_NAME := uptimer
CMD_DIR := ./cmd
GOFILES := $(shell find . -name '*.go' -not -path './vendor/*')

.PHONY: all build run test fmt vet lint docker clean

all: build

fmt:
	@echo "Formatting..."
	@gofmt -w $(GOFILES)

vet:
	@echo "Running go vet..."
	@go vet ./...

test:
	@echo "Running tests..."
	@GOCACHE=$$(mktemp -d) go test ./...

build:
	@echo "Building binary for host platform..."
	@mkdir -p bin
	@go build -o bin/$(BIN_NAME) $(CMD_DIR)

build-linux-amd64:
	@echo "Building linux/amd64 binary..."
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(BIN_NAME)-linux-amd64 $(CMD_DIR)

run:
	@echo "Running uptimer..."
	@go run $(CMD_DIR)

lint: fmt vet test

docker:
	@echo "Building docker image..."
	@docker build -t $(MOD_NAME):latest .

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin
