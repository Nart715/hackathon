NAME=component-master
BIN_NAME="component_master"
DATE_NOW=$(shell date +%Y_%m_%d_%m_%S)

# Find all .proto files
PROTO_FILES=$(shell find proto -name "*.proto")

# Main targets
.PHONY: deploy grpc clean

deploy: backup build docker_build

# Generate protobuf files
grpc:
	@echo "Generating protobuf files..."
	@for file in $(PROTO_FILES); do \
		protoc --go_out=. \
			   --go_opt=paths=source_relative \
			   --go-grpc_out=. \
			   --go-grpc_opt=paths=source_relative \
			   $$file; \
		echo "Generated $$file ✓"; \
	done

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@find proto -name "*.pb.go" -type f -delete
	@echo "Cleaned ✓"

# Backup current binary
backup:
	@echo "Backing up current binary..."
	@if [ -e bin/$(BIN_NAME) ]; then \
		mv bin/$(BIN_NAME) bin/$(BIN_NAME)_$(DATE_NOW) && \
		echo "Backup created ✓"; \
	else \
		echo "No existing binary to backup"; \
	fi

# Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	@go build -o bin/$(BIN_NAME) app/internal/server/*.go
	@echo "Build complete ✓"

# Build Docker image
docker_build:
	@echo "Building Docker image..."
	@docker build -t nartvt/$(BIN_NAME)_$(DATE_NOW) .
	@echo "Docker image built ✓"
	@docker scan --version nartvt/$(BIN_NAME)_$(DATE_NOW)
