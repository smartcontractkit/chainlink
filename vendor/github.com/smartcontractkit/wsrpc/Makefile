# Run the go linter
.PHONY: lint
lint:
	golangci-lint run

# Attempts to fix all linting issues where it can
.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: int-test-proto
int-test-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-wsrpc_out=. \
		--go-wsrpc_opt=paths=source_relative ./intgtest/internal/rpcs/rpcs.proto
