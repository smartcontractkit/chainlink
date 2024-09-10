# Name of the binary to be built
BINARY_NAME=evm-chain-bindings

# Default target
all: build

# Build the Go binary with a custom name
build:
	go build -o $(BINARY_NAME)

# Clean the build artifacts
clean:
	rm -f $(BINARY_NAME)

# Install the binary to the system's Go bin directory
install: build
	install -m 0755 $(BINARY_NAME) $(GOPATH)/bin/

# Uninstall the binary
uninstall:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Help information
help:
	@echo "Usage:"
	@echo "  make build       Build the binary with the custom name."
	@echo "  make install     Install the binary to the Go bin directory."
	@echo "  make clean       Remove the binary."
	@echo "  make uninstall   Remove the installed binary."
	@echo "  make help        Show this help message."

.PHONY: all build clean install uninstall help
