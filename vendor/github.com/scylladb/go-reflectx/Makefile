all: check test

ifndef GOBIN
export GOBIN := $(GOPATH)/bin
endif

define dl
	@curl -sSq -L $(2) -o $(GOBIN)/$(1) && chmod u+x $(GOBIN)/$(1)
endef

define dl_tgz
	@curl -sSq -L $(2) | tar zxf - --strip 1 -C $(GOBIN) --wildcards '*/$(1)'
endef

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: check
check:
	@$(GOBIN)/golangci-lint run ./...

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: bench
bench:
	@go test -tags all -run=XXX -bench=. -benchmem ./...

.PHONY: get-deps
get-deps:
	go get -t ./...

.PHONY: get-tools
get-tools:
	@echo "==> Installing tools at $(GOBIN)..."
	@$(call dl_tgz,golangci-lint,https://github.com/golangci/golangci-lint/releases/download/v1.13.1/golangci-lint-1.13.1-linux-amd64.tar.gz)
