# Chain Selectors

CCIP uses its own set of chain selectors represented by uint64 to identify blockchains. This repository contains a
mapping between the custom chain identifiers (`chainSelectorId`) chain names and the chain identifiers
used by the blockchains themselves (`chainId`).

Please refer to the [official documentation](https://docs.chain.link/ccip/supported-networks) to learn more about
supported networks and their selectors.

### Installation

`go get github.com/smartcontractkit/chain-selectors`

### Usage

```go
import (
    chainselectors "github.com/smartcontractkit/chain-selectors"
)

func main() {
    // Getting selector based on ChainId
    selector, err := chainselectors.SelectorFromChainId(420)
    
    // Getting ChainId based on ChainSelector
    chainId, err := chainselectors.ChainIdFromSelector(2664363617261496610)
    
    // Getting ChainName based on ChainId
    chainName, err := chainselectors.NameFromChainId(420)
    
    // Getting ChainId based on the ChainName
    chainId, err := chainselectors.ChainIdFromName("binance_smart_chain-testnet")
    
    // Accessing mapping directly
    lookupChainId := uint64(1337)
    if chainSelector, exists := chainselectors.EvmChainIdToChainSelector()[lookupChainId]; exists {
        fmt.Println("Found chain selector for chain", lookupChainId, ":", chainSelector)
    }
}
```

### Contributing

#### Adding new chains

Any new chains and selectors should be always added to [selectors.yml](selectors.yml) and client libraries should load
details from this file. This ensures that all client libraries are in sync and use the same mapping.
To add a new chain, please add new entry to the `selectors.yml` file and use the following format:

Make sure to run `go generate` after making any changes.

```yaml
$chain_id:
  selector: $chain_selector as uint64
  name: $chain_name as string # Although name is optional parameter, please provide it and respect the format described below
```

Chain names must respect the following format:
`<blockchain>-<type>-<network_name>-<parachain>-<rollup>-<rollup_instance>`

| Parameter | Description | Example                       |
| --- | --- |-------------------------------|
| blockchain | Name of blockchain protocol (or anchor blockchain) | `ethereum`, `cosmos`, `polkadot`    |
| type | Type of the blockchain | `testnet`, `mainnet`, `devnet`      |
| network_name | Name of specific network | `kovan`, `rinkeby`, `opera`, `kusama` |
| parachain | Name of parachain based on blockchain_protocol | `moonbeam`, `edgeware`, `okex`      |
| rollup | Name of rollup protocol | `arbitrum`, `optimism`            |
| rollup_instance | Instance of rollup | `1`, `one`                        |


[selectors.yml](selectors.yml) file is divided into sections based on the blockchain type. 
Please make sure to add new entries to the both sections and keep them sorted by chain id within these sections.

If you need to add a new chain for testing purposes (e.g. running tests with simulated environment) don't mix it with
the main file and use [test_selectors.yml](test_selectors.yml) instead. This file is used only for testing purposes.

#### Adding new client libraries

If you need a support for a new language, please open a PR with the following changes:

- Library codebase is in a separate directory
- Library uses selectors.yml as a source of truth
- Proper Github workflow is present to make sure code compiles and tests pass

