# Sample project - Basic

Most simple usage example of evm-chain-bindings tool relying on its default configuration and on a simple Solidity smart contract.

The project Makefile contains the logic to execute evm-chain-bindings script using go generate:
```
all: build

generate-abi:
	solc --abi contracts/ChainReaderTester.sol -o contracts/ --overwrite

build: generate-abi
	go generate ./...
	go build ./...
```

The file `main.go` defines the `go:generate` directive

```
//go:generate evm-chain-bindings
func main() {
}
```

This executes the `evm-chain-bindings` script using all the default options:
- Contracts ABIs will be read from contracts/ folder in the root directory
- Output will be generated in generated/evm/bindings/ folder in the root directory

To build the project and generate bindings run:
```
make
```

Then you can see the generated code in `generated/evm/bindings`