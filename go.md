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

  classDef group stroke-dasharray:6,fill:none;
  class chains,products group

	chain-selectors
	click chain-selectors href "https://github.com/smartcontractkit/chain-selectors"
	chainlink-automation --> chainlink-common
	click chainlink-automation href "https://github.com/smartcontractkit/chainlink-automation"
	chainlink-ccip --> chain-selectors
	chainlink-ccip --> chainlink-common
	click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-common --> grpc-proxy
	chainlink-common --> libocr
	click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-cosmos --> chainlink-common
	click chainlink-cosmos href "https://github.com/smartcontractkit/chainlink-cosmos"
	chainlink-data-streams --> chainlink-common
	click chainlink-data-streams href "https://github.com/smartcontractkit/chainlink-data-streams"
	chainlink-feeds --> chainlink-common
	click chainlink-feeds href "https://github.com/smartcontractkit/chainlink-feeds"
	chainlink-protos/orchestrator --> wsrpc
	click chainlink-protos/orchestrator href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-solana --> chainlink-common
	click chainlink-solana href "https://github.com/smartcontractkit/chainlink-solana"
	chainlink-starknet/relayer --> chainlink-common
	click chainlink-starknet/relayer href "https://github.com/smartcontractkit/chainlink-starknet"
	chainlink/v2 --> chainlink-automation
	chainlink/v2 --> chainlink-ccip
	chainlink/v2 --> chainlink-cosmos
	chainlink/v2 --> chainlink-data-streams
	chainlink/v2 --> chainlink-feeds
	chainlink/v2 --> chainlink-protos/orchestrator
	chainlink/v2 --> chainlink-solana
	chainlink/v2 --> chainlink-starknet/relayer
	chainlink/v2 --> tdh2/go/ocr2/decryptionplugin
	click chainlink/v2 href "https://github.com/smartcontractkit/chainlink"
	grpc-proxy
	click grpc-proxy href "https://github.com/smartcontractkit/grpc-proxy"
	libocr
	click libocr href "https://github.com/smartcontractkit/libocr"
	tdh2/go/ocr2/decryptionplugin --> libocr
	tdh2/go/ocr2/decryptionplugin --> tdh2/go/tdh2
	click tdh2/go/ocr2/decryptionplugin href "https://github.com/smartcontractkit/tdh2"
	tdh2/go/tdh2
	click tdh2/go/tdh2 href "https://github.com/smartcontractkit/tdh2"
	wsrpc
	click wsrpc href "https://github.com/smartcontractkit/wsrpc"

	subgraph tdh2-repo[tdh2]
		 tdh2/go/ocr2/decryptionplugin
		 tdh2/go/tdh2
	end
	click tdh2-repo href "https://github.com/smartcontractkit/tdh2"

	classDef outline stroke-dasharray:6,fill:none;
	class tdh2-repo outline
```
## All modules
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

  classDef group stroke-dasharray:6,fill:none;
  class chains,products group

	ccip-owner-contracts --> chain-selectors
	click ccip-owner-contracts href "https://github.com/smartcontractkit/ccip-owner-contracts"
	chain-selectors
	click chain-selectors href "https://github.com/smartcontractkit/chain-selectors"
	chainlink-automation --> chainlink-common
	click chainlink-automation href "https://github.com/smartcontractkit/chainlink-automation"
	chainlink-ccip --> chain-selectors
	chainlink-ccip --> chainlink-common
	click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-common --> grpc-proxy
	chainlink-common --> libocr
	click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-cosmos --> chainlink-common
	click chainlink-cosmos href "https://github.com/smartcontractkit/chainlink-cosmos"
	chainlink-data-streams --> chainlink-common
	click chainlink-data-streams href "https://github.com/smartcontractkit/chainlink-data-streams"
	chainlink-feeds --> chainlink-common
	click chainlink-feeds href "https://github.com/smartcontractkit/chainlink-feeds"
	chainlink-protos/job-distributor
	click chainlink-protos/job-distributor href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/orchestrator --> wsrpc
	click chainlink-protos/orchestrator href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-solana --> chainlink-common
	click chainlink-solana href "https://github.com/smartcontractkit/chainlink-solana"
	chainlink-starknet/relayer --> chainlink-common
	click chainlink-starknet/relayer href "https://github.com/smartcontractkit/chainlink-starknet"
	chainlink-testing-framework/havoc --> chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/havoc href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/lib --> chainlink-testing-framework/seth
	chainlink-testing-framework/lib --> chainlink-testing-framework/wasp
	click chainlink-testing-framework/lib href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/lib/grafana href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/seth --> seth
	click chainlink-testing-framework/seth href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/wasp --> chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/wasp href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink/core/scripts --> chainlink/deployment
	click chainlink/core/scripts href "https://github.com/smartcontractkit/chainlink"
	chainlink/deployment --> ccip-owner-contracts
	chainlink/deployment --> chainlink-protos/job-distributor
	chainlink/deployment --> chainlink-testing-framework/lib
	chainlink/deployment --> chainlink/v2
	click chainlink/deployment href "https://github.com/smartcontractkit/chainlink"
	chainlink/integration-tests --> chainlink-testing-framework/havoc
	chainlink/integration-tests --> chainlink/deployment
	click chainlink/integration-tests href "https://github.com/smartcontractkit/chainlink"
	chainlink/load-tests --> chainlink/integration-tests
	click chainlink/load-tests href "https://github.com/smartcontractkit/chainlink"
	chainlink/v2 --> chainlink-automation
	chainlink/v2 --> chainlink-ccip
	chainlink/v2 --> chainlink-cosmos
	chainlink/v2 --> chainlink-data-streams
	chainlink/v2 --> chainlink-feeds
	chainlink/v2 --> chainlink-protos/orchestrator
	chainlink/v2 --> chainlink-solana
	chainlink/v2 --> chainlink-starknet/relayer
	chainlink/v2 --> tdh2/go/ocr2/decryptionplugin
	click chainlink/v2 href "https://github.com/smartcontractkit/chainlink"
	grpc-proxy
	click grpc-proxy href "https://github.com/smartcontractkit/grpc-proxy"
	libocr
	click libocr href "https://github.com/smartcontractkit/libocr"
	seth
	click seth href "https://github.com/smartcontractkit/seth"
	tdh2/go/ocr2/decryptionplugin --> libocr
	tdh2/go/ocr2/decryptionplugin --> tdh2/go/tdh2
	click tdh2/go/ocr2/decryptionplugin href "https://github.com/smartcontractkit/tdh2"
	tdh2/go/tdh2
	click tdh2/go/tdh2 href "https://github.com/smartcontractkit/tdh2"
	wsrpc
	click wsrpc href "https://github.com/smartcontractkit/wsrpc"

	subgraph chainlink-repo[chainlink]
		 chainlink/core/scripts
		 chainlink/deployment
		 chainlink/integration-tests
		 chainlink/load-tests
		 chainlink/v2
	end
	click chainlink-repo href "https://github.com/smartcontractkit/chainlink"

	subgraph chainlink-protos-repo[chainlink-protos]
		 chainlink-protos/job-distributor
		 chainlink-protos/orchestrator
	end
	click chainlink-protos-repo href "https://github.com/smartcontractkit/chainlink-protos"

	subgraph chainlink-testing-framework-repo[chainlink-testing-framework]
		 chainlink-testing-framework/havoc
		 chainlink-testing-framework/lib
		 chainlink-testing-framework/lib/grafana
		 chainlink-testing-framework/seth
		 chainlink-testing-framework/wasp
	end
	click chainlink-testing-framework-repo href "https://github.com/smartcontractkit/chainlink-testing-framework"

	subgraph tdh2-repo[tdh2]
		 tdh2/go/ocr2/decryptionplugin
		 tdh2/go/tdh2
	end
	click tdh2-repo href "https://github.com/smartcontractkit/tdh2"

	classDef outline stroke-dasharray:6,fill:none;
	class chainlink-repo,chainlink-protos-repo,chainlink-testing-framework-repo,tdh2-repo outline
```
