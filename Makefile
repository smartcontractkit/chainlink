.PHONY: gomodtidy
gomodtidy:
	go mod tidy

.PHONY: docs
docs:
	go install golang.org/x/pkgsite/cmd/pkgsite@latest
	# http://localhost:8080/pkg/github.com/smartcontractkit/chainlink-common/pkg/
	pkgsite

PHONY: install-protoc
install-protoc:
	script/install-protoc.sh 25.1 /
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31; go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 

.PHONY: mockery
mockery: $(mockery) ## Install mockery.
	go install github.com/vektra/mockery/v2@v2.43.0

.PHONY: generate
generate: mockery install-protoc
# add our installed protoc to the head of the PATH
# maybe there is a cleaner way to do this
	 PATH=$$HOME/.local/bin:$$PATH go generate -x ./...

.PHONY: lint-workspace lint
GOLANGCI_LINT_VERSION := 1.55.2
GOLANGCI_LINT_COMMON_OPTS := --max-issues-per-linter 0 --max-same-issues 0
GOLANGCI_LINT_DIRECTORY := ./golangci-lint

lint-workspace:
	@./script/lint.sh $(GOLANGCI_LINT_VERSION) "$(GOLANGCI_LINT_COMMON_OPTS)" $(GOLANGCI_LINT_DIRECTORY)

lint:
	@./script/lint.sh $(GOLANGCI_LINT_VERSION) "$(GOLANGCI_LINT_COMMON_OPTS)" $(GOLANGCI_LINT_DIRECTORY) "--new-from-rev=origin/main"
