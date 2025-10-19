.PHONY: all build proto clean run test

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=gateway
BINARY_PATH=bin/$(BINARY_NAME)

# Protobuf
PROTO_DIR=pkg/proto
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

all: proto build

build:
	@echo "Building..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) ./cmd/gateway
	@echo "Build complete: $(BINARY_PATH)"

proto:
	@echo "Generating protobuf files..."
	@echo "Note: Requires protoc to be installed"
	@echo "Install: https://grpc.io/docs/protoc-installation/"
	@# Uncomment when protoc is available:
	@# protoc --go_out=. --go_opt=paths=source_relative \
	@#        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	@#        $(PROTO_FILES)

run: build
	@echo "Starting Voice Gateway..."
	./$(BINARY_PATH)

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f $(PROTO_DIR)/*.pb.go
	@echo "Clean complete"

test:
	$(GOTEST) -v ./...

deps:
	$(GOMOD) download
	$(GOMOD) tidy

install-tools:
	@echo "Installing development tools..."
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@echo "Tools installed. You still need to install protoc separately."
	@echo "Visit: https://grpc.io/docs/protoc-installation/"

help:
	@echo "Voice Gateway - Makefile commands:"
	@echo "  make build         - Build the gateway binary"
	@echo "  make proto         - Generate protobuf files (requires protoc)"
	@echo "  make run           - Build and run the gateway"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make install-tools - Install Go protobuf tools"
	@echo "  make help          - Show this help message"
