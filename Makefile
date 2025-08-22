# Base binary name (without extension)
BINARY_BASE := bn

# Create folder function
define MAKE_BIN_DIR
	@mkdir -p bin
endef

# Build for Windows (AMD64)
windows-amd64:
	@echo "Building for Windows with GOOS=windows, GOARCH=amd64..."
	$(call MAKE_BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows" \
		go build -o bin/$(BINARY_BASE)_windows-amd64.exe .

# Build for Windows (ARM64)
windows-arm64:
	@echo "Building for Windows (ARM) with GOOS=windows, GOARCH=arm64..."
	$(call MAKE_BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 CC="zig cc -target aarch64-windows" \
		go build -o bin/$(BINARY_BASE)-windows-arm64.exe .

# Build for Linux (AMD64)
linux-amd64:
	@echo "Building for Linux with GOOS=linux, GOARCH=amd64..."
	$(call MAKE_BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" \
		go build -o bin/$(BINARY_BASE)-linux-amd64 .

# Build for Linux (ARM64)
linux-arm64:
	@echo "Building for Linux (ARM) with GOOS=linux, GOARCH=arm64..."
	$(call MAKE_BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux" \
		go build -o bin/$(BINARY_BASE)-linux-arm64 .

# Build for macOS (ARM64)
darwin-arm64:
	@echo "Building for macOS (ARM) with GOOS=darwin, GOARCH=arm64..."
	$(call MAKE_BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 CC=clang \
		go build -o bin/$(BINARY_BASE)-macos-arm64 .

# Build all targets
all: windows-amd64 windows-arm linux-AMD64 linux-arm64 darwin-arm64
	@echo "All builds completed successfully."

# Clean target to remove generated binaries and bin folder if needed
clean:
	@echo "Cleaning generated binaries..."
	@rm -rf bin
	@rm -rf html

.PHONY: all windows-amd64 windows-arm linux-AMD64 linux-arm64 darwin-arm64 clean
