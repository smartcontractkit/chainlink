# Set default GOPATH if not already set
GOPATH ?= $(HOME)/go

.PHONY: install
install:
	go build -o $(GOPATH)/bin/local_cre -ldflags="-X main.Version=$$(git describe --tags) -X main.Commit=$$(git rev-parse HEAD) -X main.Date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" . 
	@if command -v asdf >/dev/null 2>&1; then \
		echo "Detected asdf, running reshim..."; \
		asdf reshim golang; \
	fi
	@echo ""
	@echo "✓ Local CRE has been installed successfully!"
	@echo ""
	@echo "Interactive shell is available as: local_cre sh"
	@echo ""
