# Default target
all: build

generate-abi:
	solc --abi contracts/ChainReaderTester.sol -o contracts/ --overwrite

build: generate-abi
	go generate ./...
	go build ./...