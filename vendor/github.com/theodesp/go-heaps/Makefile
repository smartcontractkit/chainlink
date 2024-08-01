
.PHONY: format
format:
	@find . -type f -name "*.go*" -print0 | xargs -0 gofmt -s -w

.PHONY: debs
debs:
	GOPATH=$(GOPATH) go get -u github.com/stretchr/testify
	GOPATH=$(GOPATH) go get -u github.com/fortytw2/leaktest

.PHONY: test
test:
	-rm coverage.txt
	@for package in $$(go list ./... | grep -v example) ; do \
		GOPATH=$(GOPATH) go test -race -coverprofile=profile.out -covermode=atomic $$package ; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt ; \
			rm profile.out ; \
		fi \
	done

.PHONY: bench
bench:
	GOPATH=$(GOPATH) go test -bench=. -check.b -benchmem

# Clean junk
.PHONY: clean
clean:
	GOPATH=$(GOPATH) go clean ./...
