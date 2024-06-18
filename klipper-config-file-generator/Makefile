# Name of the Go binary
BINARY_NAME=klipper_config_file_generator

# Go command paths (adjust if Go is not in your PATH)
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

# Directory to place the built binary
BUILD_DIR=build

# Directory containing Go source files
SRC_DIR=.

# Go source files
SRC=$(wildcard $(SRC_DIR)/*.go)

# Default target (build for Raspberry Pi)
all: build-pi

# Build the Go script for macOS
build-mac:
	@echo "Building for macOS..."
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(SRC)

# Build the Go script for Raspberry Pi (Linux ARM64)
build-pi:
	@echo "Building for Raspberry Pi (Linux ARM64)..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(SRC)

# Run the Go application on Raspberry Pi
klipper-config: build-pi
	$(BUILD_DIR)/$(BINARY_NAME)

# Run the Go application on macOS
klipper-config-mac: build-mac
	$(BUILD_DIR)/$(BINARY_NAME)

# Clean the build directory
clean:
	rm -rf $(BUILD_DIR)

# Phony targets
.PHONY: all build-mac build-pi run run-mac clean

