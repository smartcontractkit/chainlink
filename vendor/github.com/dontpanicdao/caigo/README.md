<h1 align="center">Golang Library for StarkNet</h1>

<p align="center">
    <a href="https://pkg.go.dev/github.com/dontpanicdao/caigo">
        <img src="https://pkg.go.dev/badge/github.com/dontpanicdao/caigo.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/dontpanicdao/caigo/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-black">
    </a>
    <a href="https://starkware.co/">
        <img src="https://img.shields.io/badge/powered_by-StarkWare-navy">
    </a>
    <a href="https://github.com/dontpanicdao/caigo/actions/workflows/go.yml">
        <img src="https://github.com/dontpanicdao/caigo/actions/workflows/go.yml/badge.svg?branch=main" alt="Go">
    </a>
</p>

Caigo is an MIT-licensed Go library for interacting with [StarkNet](https://docs.starknet.io/docs/intro).

### Getting Started

- library documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/dontpanicdao/caigo).
- [curve example](./examples/curve) initializing the StarkCurve for signing and verification
- [contract example](./examples/contract) for smart contract deployment and function call
- [account example](./examples/contract) for Account initialization and invocation call

### Run Examples

***starkcurve***

```sh
cd examples/curve
go mod tidy
go run main.go
```

***starknet contract***

```sh
cd examples/contract
go mod tidy
go run main.go
```

***starknet account***

```sh
cd examples/account
go mod tidy
go run main.go
```

### RPC

Caigo RPC implements the [StarkNet RPC Spec](https://github.com/starkware-libs/starknet-specs):

| Method                                         | Implemented (*)    |
| ---------------------------------------------- | ------------------ |
| `starknet_getBlockByHash`                      | :heavy_check_mark: |
| `starknet_getBlockByNumber`                    | :heavy_check_mark: |
| `starknet_getTransactionByHash`                | :heavy_check_mark: |
| `starknet_getTransactionReceipt`               | :heavy_check_mark: |
| `starknet_getClass`                            | :heavy_check_mark: |
| `starknet_getClassHashAt`                      | :heavy_check_mark: |
| `starknet_getClassAt`                          | :heavy_check_mark: |
| `starknet_call`                                | :heavy_check_mark: |
| `starknet_blockNumber`                         | :heavy_check_mark: |
| `starknet_chainId`                             | :heavy_check_mark: |
| `starknet_syncing`                             | :heavy_check_mark: |
| `starknet_getEvents`                           | :heavy_check_mark: |
| `starknet_addInvokeTransaction`                | :heavy_check_mark: |
| `starknet_addDeployTransaction`                | :heavy_check_mark: |
| `starknet_addDeclareTransaction`               | :heavy_check_mark: |
| `starknet_estimateFee`                         | :heavy_check_mark: |
| `starknet_getBlockTransactionCountByHash`      | :heavy_check_mark: |
| `starknet_getBlockTransactionCountByNumber`    | :heavy_check_mark: |
| `starknet_getTransactionByBlockNumberAndIndex` | :heavy_check_mark: |
| `starknet_getTransactionByBlockHashAndIndex`   | :heavy_check_mark: |
| `starknet_getStorageAt`                        | :heavy_check_mark: |
| `starknet_getNonce`                            | :heavy_check_mark: |
| `starknet_getStateUpdate`                      | :heavy_check_mark: |
| *`starknet_traceBlockTransactions`             | :x:                |
| *`starknet_traceTransaction`                   | :x:                |

> (*) some methods are not implemented because they are not yet available
> from [eqlabs/pathfinder](https://github.com/eqlabs/pathfinder).

### Run Tests

```go
go test -v ./...
```

### Run RPC Tests

```go
go test -v ./rpc -env [mainnet|devnet|testnet|mock]
```

### Run Benchmarks

```go
go test -bench=.
```

### Compatibility and stability

Caigo is currently under active development and will under go breaking changes until the initial stable(v1.0.0) release. The example directories and *_test.go files should always be applicable for the latest commitment on the main branch.
*NOTE: examples and tests may be out of sync with tagged versions and pkg.go.dev documentation*

## Issues

If you find an issue/bug or have a feature request please submit an issue here
[Issues](https://github.com/dontpanicdao/caigo/issues)

## Contributing

If you are looking to contribute, please head to the
[Contributing](https://github.com/dontpanicdao/caigo/blob/main/CONTRIBUTING.md) section.
