.PHONY: build
build:
	go build -o observability-lib .

.PHONY: run
run:
	go run main.go

.PHONY: lint
lint:
	golangci-lint run

test:
	go test -v ./...