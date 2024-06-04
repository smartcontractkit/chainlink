
ensure_go_1_21:
	@go version | grep -q 'go1.21' || (echo "Please use go1.21" && exit 1)

test: ensure_go_1_21
	go test -race -fullpath -shuffle on -count 10 ./...
