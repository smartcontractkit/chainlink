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
	go install github.com/vektra/mockery/v2@v2.38.0

.PHONY: generate
generate: mockery install-protoc
# add our installed protoc to the head of the PATH
# maybe there is a cleaner way to do this
	 PATH=$$HOME/.local/bin:$$PATH go generate -x ./...

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint for all issues.
	[ -d "./golangci-lint" ] || mkdir ./golangci-lint && \
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.55.2 golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 > ./golangci-lint/$(shell date +%Y-%m-%d_%H:%M:%S).txt
