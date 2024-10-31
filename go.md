# smartcontractkit Go modules
## Main module
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

  subgraph chainlink-protos
    chainlink-protos/orchestrator
    chainlink-protos/job-distributor
  end

  classDef outline stroke-dasharray:6,fill:none;
  class chains,products,tdh2,chainlink-protos outline

  chainlink/v2 --> chain-selectors
  click chain-selectors href "https://github.com/smartcontractkit/chain-selectors"
  chainlink/v2 --> chainlink-automation
  click chainlink-automation href "https://github.com/smartcontractkit/chainlink-automation"
  chainlink/v2 --> chainlink-ccip
  click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
  chainlink/v2 --> chainlink-common
  click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
  chainlink/v2 --> chainlink-cosmos
  click chainlink-cosmos href "https://github.com/smartcontractkit/chainlink-cosmos"
  chainlink/v2 --> chainlink-data-streams
  click chainlink-data-streams href "https://github.com/smartcontractkit/chainlink-data-streams"
  chainlink/v2 --> chainlink-feeds
  click chainlink-feeds href "https://github.com/smartcontractkit/chainlink-feeds"
  chainlink/v2 --> chainlink-protos/orchestrator
  click chainlink-protos/orchestrator href "https://github.com/smartcontractkit/chainlink-protos"
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
  chainlink-ccip --> chain-selectors
  chainlink-ccip --> chainlink-common
  chainlink-ccip --> libocr
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
  chainlink-protos/orchestrator --> wsrpc
  chainlink-solana --> chainlink-common
  chainlink-solana --> libocr
  chainlink-solana --> grpc-proxy
  chainlink-starknet/relayer --> chainlink-common
  chainlink-starknet/relayer --> libocr
  chainlink-starknet/relayer --> grpc-proxy
  tdh2/go/ocr2/decryptionplugin --> libocr
  tdh2/go/ocr2/decryptionplugin --> tdh2/go/tdh2
```
## All modules
```mermaid
flowchart LR
  subgraph chainlink
    chainlink/v2
    chainlink/integration-tests
    chainlink/load-tests
    chainlink/core/scripts
  end

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

  subgraph chainlink-testing-framework
    chainlink-testing-framework/grafana
    chainlink-testing-framework/havoc
    chainlink-testing-framework/lib
    chainlink-testing-framework/lib/grafana
    chainlink-testing-framework/seth
    chainlink-testing-framework/wasp
  end

  subgraph chainlink-protos
    chainlink-protos/orchestrator
    chainlink-protos/job-distributor
  end

  classDef outline stroke-dasharray:6,fill:none;
  class chainlink,chains,products,tdh2,chainlink-protos,chainlink-testing-framework outline

  	chainlink/v2 --> chain-selectors
  click chain-selectors href "https://github.com/smartcontractkit/chain-selectors"
  	chainlink/v2 --> chainlink-automation
  click chainlink-automation href "https://github.com/smartcontractkit/chainlink-automation"
  	chainlink/v2 --> chainlink-ccip
  click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
  	chainlink/v2 --> chainlink-common
  click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
  	chainlink/v2 --> chainlink-cosmos
  click chainlink-cosmos href "https://github.com/smartcontractkit/chainlink-cosmos"
  	chainlink/v2 --> chainlink-data-streams
  click chainlink-data-streams href "https://github.com/smartcontractkit/chainlink-data-streams"
  	chainlink/v2 --> chainlink-feeds
  click chainlink-feeds href "https://github.com/smartcontractkit/chainlink-feeds"
  	chainlink/v2 --> chainlink-protos/orchestrator
  click chainlink-protos/orchestrator href "https://github.com/smartcontractkit/chainlink-protos"
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
  	chainlink-ccip --> chain-selectors
  	chainlink-ccip --> chainlink-common
  	chainlink-ccip --> libocr
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
  	chainlink-protos/orchestrator --> wsrpc
  	chainlink-solana --> chainlink-common
  	chainlink-solana --> libocr
  	chainlink-solana --> grpc-proxy
  	chainlink-starknet/relayer --> chainlink-common
  	chainlink-starknet/relayer --> libocr
  	chainlink-starknet/relayer --> grpc-proxy
  	tdh2/go/ocr2/decryptionplugin --> libocr
  	tdh2/go/ocr2/decryptionplugin --> tdh2/go/tdh2
  	chainlink/core/scripts --> ccip-owner-contracts
  click ccip-owner-contracts href "https://github.com/smartcontractkit/ccip-owner-contracts"
  	chainlink/core/scripts --> chain-selectors
  	chainlink/core/scripts --> chainlink-automation
  	chainlink/core/scripts --> chainlink-ccip
  	chainlink/core/scripts --> chainlink-common
  	chainlink/core/scripts --> chainlink-cosmos
  	chainlink/core/scripts --> chainlink-data-streams
  	chainlink/core/scripts --> chainlink-feeds
  	chainlink/core/scripts --> chainlink-protos/job-distributor
  click chainlink-protos/job-distributor href "https://github.com/smartcontractkit/chainlink-protos"
  	chainlink/core/scripts --> chainlink-protos/orchestrator
  	chainlink/core/scripts --> chainlink-solana
  	chainlink/core/scripts --> chainlink-starknet/relayer
  	chainlink/core/scripts --> chainlink/deployment
  click chainlink/deployment href "https://github.com/smartcontractkit/chainlink"
  	chainlink/core/scripts --> chainlink/v2
  click chainlink/v2 href "https://github.com/smartcontractkit/chainlink"
  	chainlink/core/scripts --> grpc-proxy
  	chainlink/core/scripts --> libocr
  	chainlink/core/scripts --> tdh2/go/ocr2/decryptionplugin
  	chainlink/core/scripts --> tdh2/go/tdh2
  	chainlink/core/scripts --> wsrpc
  	ccip-owner-contracts --> chain-selectors
  	chainlink/deployment --> ccip-owner-contracts
  	chainlink/deployment --> chain-selectors
  	chainlink/deployment --> chainlink-ccip
  	chainlink/deployment --> chainlink-common
  	chainlink/deployment --> chainlink-protos/job-distributor
  	chainlink/deployment --> chainlink-testing-framework/lib
  click chainlink-testing-framework/lib href "https://github.com/smartcontractkit/chainlink-testing-framework"
  	chainlink/deployment --> chainlink/v2
  	chainlink/deployment --> libocr
  	chainlink/deployment --> chainlink-automation
  	chainlink/deployment --> chainlink-cosmos
  	chainlink/deployment --> chainlink-data-streams
  	chainlink/deployment --> chainlink-feeds
  	chainlink/deployment --> chainlink-protos/orchestrator
  	chainlink/deployment --> chainlink-solana
  	chainlink/deployment --> chainlink-starknet/relayer
  	chainlink/deployment --> chainlink-testing-framework/grafana
  click chainlink-testing-framework/grafana href "https://github.com/smartcontractkit/chainlink-testing-framework"
  	chainlink/deployment --> chainlink-testing-framework/seth
  click chainlink-testing-framework/seth href "https://github.com/smartcontractkit/chainlink-testing-framework"
  	chainlink/deployment --> chainlink-testing-framework/wasp
  click chainlink-testing-framework/wasp href "https://github.com/smartcontractkit/chainlink-testing-framework"
  	chainlink/deployment --> grpc-proxy
  	chainlink/deployment --> tdh2/go/ocr2/decryptionplugin
  	chainlink/deployment --> tdh2/go/tdh2
  	chainlink/deployment --> wsrpc
  	chainlink-testing-framework/lib --> chainlink-testing-framework/seth
  	chainlink-testing-framework/lib --> chainlink-testing-framework/wasp
  	chainlink-testing-framework/lib --> chainlink-testing-framework/grafana
  	chainlink-testing-framework/seth --> seth
  click seth href "https://github.com/smartcontractkit/seth"
  	chainlink-testing-framework/wasp --> chainlink-testing-framework/grafana
  	chainlink/integration-tests --> ccip-owner-contracts
  	chainlink/integration-tests --> chain-selectors
  	chainlink/integration-tests --> chainlink-automation
  	chainlink/integration-tests --> chainlink-ccip
  	chainlink/integration-tests --> chainlink-common
  	chainlink/integration-tests --> chainlink-cosmos
  	chainlink/integration-tests --> chainlink-data-streams
  	chainlink/integration-tests --> chainlink-feeds
  	chainlink/integration-tests --> chainlink-protos/job-distributor
  	chainlink/integration-tests --> chainlink-protos/orchestrator
  	chainlink/integration-tests --> chainlink-solana
  	chainlink/integration-tests --> chainlink-starknet/relayer
  	chainlink/integration-tests --> chainlink-testing-framework/havoc
  click chainlink-testing-framework/havoc href "https://github.com/smartcontractkit/chainlink-testing-framework"
  	chainlink/integration-tests --> chainlink-testing-framework/lib
  	chainlink/integration-tests --> chainlink-testing-framework/lib/grafana
  click chainlink-testing-framework/lib/grafana href "https://github.com/smartcontractkit/chainlink-testing-framework"
  	chainlink/integration-tests --> chainlink-testing-framework/seth
  	chainlink/integration-tests --> chainlink-testing-framework/wasp
  	chainlink/integration-tests --> chainlink/deployment
  	chainlink/integration-tests --> chainlink/v2
  	chainlink/integration-tests --> grpc-proxy
  	chainlink/integration-tests --> libocr
  	chainlink/integration-tests --> tdh2/go/ocr2/decryptionplugin
  	chainlink/integration-tests --> tdh2/go/tdh2
  	chainlink/integration-tests --> wsrpc
  	chainlink-testing-framework/havoc --> chainlink-testing-framework/lib/grafana
  	chainlink-testing-framework/lib --> chainlink-testing-framework/lib/grafana
  	chainlink-testing-framework/wasp --> chainlink-testing-framework/lib/grafana
  	chainlink/load-tests --> chain-selectors
  	chainlink/load-tests --> chainlink-automation
  	chainlink/load-tests --> chainlink-ccip
  	chainlink/load-tests --> chainlink-common
  	chainlink/load-tests --> chainlink-cosmos
  	chainlink/load-tests --> chainlink-data-streams
  	chainlink/load-tests --> chainlink-feeds
  	chainlink/load-tests --> chainlink-protos/orchestrator
  	chainlink/load-tests --> chainlink-solana
  	chainlink/load-tests --> chainlink-starknet/relayer
  	chainlink/load-tests --> chainlink-testing-framework/havoc
  	chainlink/load-tests --> chainlink-testing-framework/lib
  	chainlink/load-tests --> chainlink-testing-framework/lib/grafana
  	chainlink/load-tests --> chainlink-testing-framework/seth
  	chainlink/load-tests --> chainlink-testing-framework/wasp
  	chainlink/load-tests --> chainlink/deployment
  	chainlink/load-tests --> chainlink/integration-tests
  click chainlink/integration-tests href "https://github.com/smartcontractkit/chainlink"
  	chainlink/load-tests --> chainlink/v2
  	chainlink/load-tests --> grpc-proxy
  	chainlink/load-tests --> libocr
  	chainlink/load-tests --> tdh2/go/ocr2/decryptionplugin
  	chainlink/load-tests --> tdh2/go/tdh2
  	chainlink/load-tests --> wsrpc
```
