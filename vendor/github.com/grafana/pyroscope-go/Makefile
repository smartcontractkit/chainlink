GO_VERSION_PRE20 := $(shell go version  | awk '{print $$3}' | awk -F '.' '{print ($$1 == "go1" && int($$2) < 20)}')
TEST_PACKAGES := ./... ./godeltaprof/compat/... ./godeltaprof/...

.PHONY: test
test:
	go test -race $(shell go list $(TEST_PACKAGES) | grep -v /example)

.PHONY: go/mod
go/mod:
	GO111MODULE=on go mod download
	go work sync
	GO111MODULE=on go mod tidy
	cd godeltaprof/compat/ && GO111MODULE=on go mod download
	cd godeltaprof/compat/ && GO111MODULE=on go mod tidy
	cd godeltaprof/  && GO111MODULE=on go mod download
	cd godeltaprof/ && GO111MODULE=on go mod tidy

.PHONY: go/mod_16_for_testing
go/mod_16_for_testing:
	rm -rf godeltaprof/compat/go.mod godeltaprof/compat/go.sum godeltaprof/go.mod godeltaprof/go.sum go.work otelpyroscope/
	cat go.mod_go16_test.txt > go.mod
	go mod tidy
