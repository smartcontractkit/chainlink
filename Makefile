.DEFAULT_GOAL := build
.PHONY: build

LDFLAGS=-ldflags "-X github.com/smartcontractkit/chainlink/store.Sha=`git rev-parse HEAD`"

build:
	@dep ensure
	@go build $(LDFLAGS) -o chainlink
