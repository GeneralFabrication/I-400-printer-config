# Name of the Go binary
BINARY_NAME=config_file_generator

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

# Run the addconfig command on Raspberry Pi
addconfig: build-pi
	$(BUILD_DIR)/$(BINARY_NAME) addconfig

# Run the updateconfig command on Raspberry Pi
updateconfig: build-pi
	$(BUILD_DIR)/$(BINARY_NAME) updateconfig

# Run the rollback command on Raspberry Pi
rollback: build-pi
	$(BUILD_DIR)/$(BINARY_NAME) rollback

# Run the addconfig command on macOS
addconfig-mac: build-mac
	$(BUILD_DIR)/$(BINARY_NAME) addconfig

# Run the updateconfig command on macOS
updateconfig-mac: build-mac
	$(BUILD_DIR)/$(BINARY_NAME) updateconfig

# Run the rollback command on macOS
rollback-mac: build-mac
	$(BUILD_DIR)/$(BINARY_NAME) rollback

# Clean the build directory
clean:
	rm -rf $(BUILD_DIR)

# Phony targets
.PHONY: all build-mac build-pi addconfig updateconfig rollback addconfig-mac updateconfig-mac rollback-mac clean

