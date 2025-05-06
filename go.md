# smartcontractkit Go modules
## Main module
```mermaid
flowchart LR
  subgraph chains
    chainlink-aptos
    chainlink-cosmos
    chainlink-evm
    chainlink-solana
    chainlink-starknet/relayer
    chainlink-tron/relayer
  end

  subgraph products
    chainlink-automation
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
	chainlink-ccip --> chainlink-protos/rmn/v1.6/go
	click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-ccip/chains/solana --> chainlink-ccip
	click chainlink-ccip/chains/solana href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-common --> freeport
	chainlink-common --> grpc-proxy
	chainlink-common --> libocr
	click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-common/pkg/monitoring
	click chainlink-common/pkg/monitoring href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-data-streams --> chainlink-common
	click chainlink-data-streams href "https://github.com/smartcontractkit/chainlink-data-streams"
	chainlink-evm --> chainlink-framework/capabilities
	chainlink-evm --> chainlink-framework/chains
	chainlink-evm --> chainlink-protos/svr
	click chainlink-evm href "https://github.com/smartcontractkit/chainlink-evm"
	chainlink-feeds --> chainlink-common
	click chainlink-feeds href "https://github.com/smartcontractkit/chainlink-feeds"
	chainlink-framework/capabilities
	click chainlink-framework/capabilities href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-framework/chains --> chainlink-framework/multinode
	click chainlink-framework/chains href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-framework/metrics --> chainlink-common
	click chainlink-framework/metrics href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-framework/multinode --> chainlink-framework/metrics
	click chainlink-framework/multinode href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-protos/billing/go --> chainlink-protos/workflows/go
	click chainlink-protos/billing/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/orchestrator --> wsrpc
	click chainlink-protos/orchestrator href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/rmn/v1.6/go
	click chainlink-protos/rmn/v1.6/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/svr
	click chainlink-protos/svr href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/workflows/go
	click chainlink-protos/workflows/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-solana --> chainlink-ccip/chains/solana
	chainlink-solana --> chainlink-common/pkg/monitoring
	chainlink-solana --> chainlink-framework/multinode
	click chainlink-solana href "https://github.com/smartcontractkit/chainlink-solana"
	chainlink-tron/relayer --> chainlink-evm
	click chainlink-tron/relayer href "https://github.com/smartcontractkit/chainlink-tron"
	chainlink/v2 --> chainlink-automation
	chainlink/v2 --> chainlink-data-streams
	chainlink/v2 --> chainlink-feeds
	chainlink/v2 --> chainlink-protos/billing/go
	chainlink/v2 --> chainlink-protos/orchestrator
	chainlink/v2 --> chainlink-solana
	chainlink/v2 --> chainlink-tron/relayer
	chainlink/v2 --> tdh2/go/ocr2/decryptionplugin
	click chainlink/v2 href "https://github.com/smartcontractkit/chainlink"
	freeport
	click freeport href "https://github.com/smartcontractkit/freeport"
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

	subgraph chainlink-ccip-repo[chainlink-ccip]
		 chainlink-ccip
		 chainlink-ccip/chains/solana
	end
	click chainlink-ccip-repo href "https://github.com/smartcontractkit/chainlink-ccip"

	subgraph chainlink-common-repo[chainlink-common]
		 chainlink-common
		 chainlink-common/pkg/monitoring
	end
	click chainlink-common-repo href "https://github.com/smartcontractkit/chainlink-common"

	subgraph chainlink-framework-repo[chainlink-framework]
		 chainlink-framework/capabilities
		 chainlink-framework/chains
		 chainlink-framework/metrics
		 chainlink-framework/multinode
	end
	click chainlink-framework-repo href "https://github.com/smartcontractkit/chainlink-framework"

	subgraph chainlink-protos-repo[chainlink-protos]
		 chainlink-protos/billing/go
		 chainlink-protos/orchestrator
		 chainlink-protos/rmn/v1.6/go
		 chainlink-protos/svr
		 chainlink-protos/workflows/go
	end
	click chainlink-protos-repo href "https://github.com/smartcontractkit/chainlink-protos"

	subgraph tdh2-repo[tdh2]
		 tdh2/go/ocr2/decryptionplugin
		 tdh2/go/tdh2
	end
	click tdh2-repo href "https://github.com/smartcontractkit/tdh2"

	classDef outline stroke-dasharray:6,fill:none;
	class chainlink-ccip-repo,chainlink-common-repo,chainlink-framework-repo,chainlink-protos-repo,tdh2-repo outline
```
## All modules
```mermaid
flowchart LR
  subgraph chains
    chainlink-aptos
    chainlink-cosmos
    chainlink-evm
    chainlink-solana
    chainlink-starknet/relayer
    chainlink-tron/relayer
  end

  subgraph products
    chainlink-automation
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
	chainlink-aptos --> chainlink-common
	click chainlink-aptos href "https://github.com/smartcontractkit/chainlink-aptos"
	chainlink-automation --> chainlink-common
	click chainlink-automation href "https://github.com/smartcontractkit/chainlink-automation"
	chainlink-ccip --> chain-selectors
	chainlink-ccip --> chainlink-common
	chainlink-ccip --> chainlink-protos/rmn/v1.6/go
	click chainlink-ccip href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-ccip/chains/solana --> chainlink-ccip
	click chainlink-ccip/chains/solana href "https://github.com/smartcontractkit/chainlink-ccip"
	chainlink-common --> freeport
	chainlink-common --> grpc-proxy
	chainlink-common --> libocr
	click chainlink-common href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-common/pkg/monitoring
	click chainlink-common/pkg/monitoring href "https://github.com/smartcontractkit/chainlink-common"
	chainlink-data-streams --> chainlink-common
	click chainlink-data-streams href "https://github.com/smartcontractkit/chainlink-data-streams"
	chainlink-deployments-framework --> ccip-owner-contracts
	chainlink-deployments-framework --> chainlink-protos/job-distributor
	chainlink-deployments-framework --> mcms
	click chainlink-deployments-framework href "https://github.com/smartcontractkit/chainlink-deployments-framework"
	chainlink-evm --> chainlink-framework/capabilities
	chainlink-evm --> chainlink-framework/chains
	chainlink-evm --> chainlink-protos/svr
	click chainlink-evm href "https://github.com/smartcontractkit/chainlink-evm"
	chainlink-feeds --> chainlink-common
	click chainlink-feeds href "https://github.com/smartcontractkit/chainlink-feeds"
	chainlink-framework/capabilities
	click chainlink-framework/capabilities href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-framework/chains --> chainlink-framework/multinode
	click chainlink-framework/chains href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-framework/metrics --> chainlink-common
	click chainlink-framework/metrics href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-framework/multinode --> chainlink-framework/metrics
	click chainlink-framework/multinode href "https://github.com/smartcontractkit/chainlink-framework"
	chainlink-protos/billing/go --> chainlink-protos/workflows/go
	click chainlink-protos/billing/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/job-distributor
	click chainlink-protos/job-distributor href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/orchestrator --> wsrpc
	click chainlink-protos/orchestrator href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/rmn/v1.6/go
	click chainlink-protos/rmn/v1.6/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/svr
	click chainlink-protos/svr href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-protos/workflows/go
	click chainlink-protos/workflows/go href "https://github.com/smartcontractkit/chainlink-protos"
	chainlink-solana --> chainlink-ccip/chains/solana
	chainlink-solana --> chainlink-common/pkg/monitoring
	chainlink-solana --> chainlink-framework/multinode
	click chainlink-solana href "https://github.com/smartcontractkit/chainlink-solana"
	chainlink-testing-framework/framework
	click chainlink-testing-framework/framework href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/havoc --> chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/havoc href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/lib --> chainlink-testing-framework/parrot
	chainlink-testing-framework/lib --> chainlink-testing-framework/seth
	click chainlink-testing-framework/lib href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/lib/grafana href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/parrot
	click chainlink-testing-framework/parrot href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/sentinel --> chainlink-testing-framework/lib
	click chainlink-testing-framework/sentinel href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/seth
	click chainlink-testing-framework/seth href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/wasp --> chainlink-testing-framework/lib
	chainlink-testing-framework/wasp --> chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/wasp href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-tron/relayer --> chainlink-evm
	click chainlink-tron/relayer href "https://github.com/smartcontractkit/chainlink-tron"
	chainlink/core/scripts --> chainlink/deployment
	click chainlink/core/scripts href "https://github.com/smartcontractkit/chainlink"
	chainlink/deployment --> chainlink-deployments-framework
	chainlink/deployment --> chainlink-testing-framework/lib
	chainlink/deployment --> chainlink/v2
	click chainlink/deployment href "https://github.com/smartcontractkit/chainlink"
	chainlink/integration-tests --> chainlink-testing-framework/havoc
	chainlink/integration-tests --> chainlink-testing-framework/sentinel
	chainlink/integration-tests --> chainlink-testing-framework/wasp
	chainlink/integration-tests --> chainlink/deployment
	click chainlink/integration-tests href "https://github.com/smartcontractkit/chainlink"
	chainlink/load-tests --> chainlink/integration-tests
	click chainlink/load-tests href "https://github.com/smartcontractkit/chainlink"
	chainlink/system-tests/lib --> chainlink/deployment
	click chainlink/system-tests/lib href "https://github.com/smartcontractkit/chainlink"
	chainlink/system-tests/tests --> chainlink-testing-framework/havoc
	chainlink/system-tests/tests --> chainlink-testing-framework/wasp
	chainlink/system-tests/tests --> chainlink/system-tests/lib
	click chainlink/system-tests/tests href "https://github.com/smartcontractkit/chainlink"
	chainlink/v2 --> chainlink-automation
	chainlink/v2 --> chainlink-data-streams
	chainlink/v2 --> chainlink-feeds
	chainlink/v2 --> chainlink-protos/billing/go
	chainlink/v2 --> chainlink-protos/orchestrator
	chainlink/v2 --> chainlink-solana
	chainlink/v2 --> chainlink-tron/relayer
	chainlink/v2 --> tdh2/go/ocr2/decryptionplugin
	click chainlink/v2 href "https://github.com/smartcontractkit/chainlink"
	freeport
	click freeport href "https://github.com/smartcontractkit/freeport"
	grpc-proxy
	click grpc-proxy href "https://github.com/smartcontractkit/grpc-proxy"
	libocr
	click libocr href "https://github.com/smartcontractkit/libocr"
	mcms --> chainlink-aptos
	mcms --> chainlink-ccip/chains/solana
	mcms --> chainlink-testing-framework/framework
	click mcms href "https://github.com/smartcontractkit/mcms"
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
		 chainlink/system-tests/lib
		 chainlink/system-tests/tests
		 chainlink/v2
	end
	click chainlink-repo href "https://github.com/smartcontractkit/chainlink"

	subgraph chainlink-ccip-repo[chainlink-ccip]
		 chainlink-ccip
		 chainlink-ccip/chains/solana
	end
	click chainlink-ccip-repo href "https://github.com/smartcontractkit/chainlink-ccip"

	subgraph chainlink-common-repo[chainlink-common]
		 chainlink-common
		 chainlink-common/pkg/monitoring
	end
	click chainlink-common-repo href "https://github.com/smartcontractkit/chainlink-common"

	subgraph chainlink-framework-repo[chainlink-framework]
		 chainlink-framework/capabilities
		 chainlink-framework/chains
		 chainlink-framework/metrics
		 chainlink-framework/multinode
	end
	click chainlink-framework-repo href "https://github.com/smartcontractkit/chainlink-framework"

	subgraph chainlink-protos-repo[chainlink-protos]
		 chainlink-protos/billing/go
		 chainlink-protos/job-distributor
		 chainlink-protos/orchestrator
		 chainlink-protos/rmn/v1.6/go
		 chainlink-protos/svr
		 chainlink-protos/workflows/go
	end
	click chainlink-protos-repo href "https://github.com/smartcontractkit/chainlink-protos"

	subgraph chainlink-testing-framework-repo[chainlink-testing-framework]
		 chainlink-testing-framework/framework
		 chainlink-testing-framework/havoc
		 chainlink-testing-framework/lib
		 chainlink-testing-framework/lib/grafana
		 chainlink-testing-framework/parrot
		 chainlink-testing-framework/sentinel
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
	class chainlink-repo,chainlink-ccip-repo,chainlink-common-repo,chainlink-framework-repo,chainlink-protos-repo,chainlink-testing-framework-repo,tdh2-repo outline
```
