## RPC implementation

starknet.go RPC implementation provides the RPC API to perform operations with 
Starknet. It is currently being tested and maintained up-to-date with
Pathfinder and relies on [go-ethereum](github.com/ethereum/go-ethereum/rpc)
to provide the JSON RPC 2.0 client implementation.

If you need starknet.go to support another API, open an issue on the project.

### Testing the RPC API

To test the RPC API, you should simply go the the rpc directory and run
`go test` like below:

```shell
cd rpc
go test -v .
```

We provide an additional `-env` flag to `go test` so that you can choose the
environment you want to test. For instance, if you plan to test with the
`testnet`, run:

```shell
cd rpc
go test -env testnet -v .
```

Supported environments are `mock`, `testnet` and `mainnet`. The support for
`devnet` is planned but might require some dedicated condition since it is empty. 

If you plan to specify an alternative URL to test the environment, you can set
the `INTEGRATION_BASE` environment variable. In addition, tests load `.env.${env}`,
and `.env` before relying on the environment variable. So for instanve if you want
the URL to change only for the testnet environment, you could add the line below
in `.env.testnet`:

```text
INTEGRATION_BASE=http://localhost:9546
```
