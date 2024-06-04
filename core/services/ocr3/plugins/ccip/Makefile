
ensure_go_1_21:
	@go version | grep -q 'go1.21' || (echo "Please use go1.21" && exit 1)

ensure_golangcilint_1_59:
	@golangci-lint --version | grep -q '1.59' || (echo "Please use golangci-lint 1.59" && exit 1)

test: ensure_go_1_21
	go test -race -fullpath -shuffle on -count 10 ./...

lint: ensure_go_1_21
	golangci-lint run
