# smartcontractkit Go modules
```mermaid
flowchart LR
  subgraph chains
    chainlink-cosmos
    chainlink-solana
    chainlink-starknet/relayer
    chainlink-evm
  end

  subgraph products
    chainlink-automation
    chainlink-ccip
    chainlink-data-streams
    chainlink-feeds
    chainlink-functions
    chainlink-vrf
  end

  subgraph tdh2
    tdh2/go/tdh2
    tdh2/go/ocr2/decryptionplugin
  end

  classDef outline stroke-dasharray:6,fill:none;
  class chains,products,tdh2 outline

  chainlink/v2 --> chain-selectors
  click chain-selectors href "https://github.com/smartcontractkit/chain-selectors"
  chainlink/v2 --> chainlink-automation
  click chainlink-automation href "https://github.com/smartcontractkit/chainlink-automation"
  chainlink/v2 --> chainlink-common
  click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
  chainlink/v2 --> chainlink-cosmos
  click chainlink-cosmos href "https://github.com/smartcontractkit/chainlink-cosmos"
  chainlink/v2 --> chainlink-data-streams
  click chainlink-data-streams href "https://github.com/smartcontractkit/chainlink-data-streams"
  chainlink/v2 --> chainlink-feeds
  click chainlink-feeds href "https://github.com/smartcontractkit/chainlink-feeds"
  chainlink/v2 --> chainlink-solana
  click chainlink-solana href "https://github.com/smartcontractkit/chainlink-solana"
  chainlink/v2 --> chainlink-starknet/relayer
  click chainlink-starknet/relayer href "https://github.com/smartcontractkit/chainlink-starknet"
  chainlink/v2 --> grpc-proxy
  click grpc-proxy href "https://github.com/smartcontractkit/grpc-proxy"
  chainlink/v2 --> libocr
  click libocr href "https://github.com/smartcontractkit/libocr"
  chainlink/v2 --> tdh2/go/ocr2/decryptionplugin
  click tdh2/go/ocr2/decryptionplugin href "https://github.com/smartcontractkit/tdh2"
  chainlink/v2 --> tdh2/go/tdh2
  click tdh2/go/tdh2 href "https://github.com/smartcontractkit/tdh2"
  chainlink/v2 --> wsrpc
  click wsrpc href "https://github.com/smartcontractkit/wsrpc"
  chainlink-automation --> chainlink-common
  chainlink-automation --> libocr
  chainlink-common --> grpc-proxy
  chainlink-common --> libocr
  chainlink-cosmos --> chainlink-common
  chainlink-cosmos --> libocr
  chainlink-cosmos --> grpc-proxy
  chainlink-data-streams --> chainlink-common
  chainlink-data-streams --> libocr
  chainlink-data-streams --> grpc-proxy
  chainlink-feeds --> chainlink-common
  chainlink-feeds --> libocr
  chainlink-feeds --> grpc-proxy
  chainlink-solana --> chainlink-common
  chainlink-solana --> libocr
  chainlink-solana --> grpc-proxy
  chainlink-starknet/relayer --> chainlink-common
  chainlink-starknet/relayer --> libocr
  chainlink-starknet/relayer --> grpc-proxy
  tdh2/go/ocr2/decryptionplugin --> libocr
  tdh2/go/ocr2/decryptionplugin --> tdh2/go/tdh2
```
