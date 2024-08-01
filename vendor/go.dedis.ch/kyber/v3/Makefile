tidy:
	#go install honnef.co/go/tools/cmd/staticcheck@latest
	go mod tidy

generate: tidy
	go generate ./...

# Coding style static check.
lint: tidy
	@echo "Please setup a linter!"
	#golangci-lint run
	#staticcheck go list ./...


vet: tidy
	go vet ./...

test: tidy
	go test ./...

coverage: tidy
	go test -json -covermode=count -coverprofile=profile.cov ./... > report.json

# target to run all the possible checks; it's a good habit to run it before
# pushing code
check: lint vet test
	echo "check done"
