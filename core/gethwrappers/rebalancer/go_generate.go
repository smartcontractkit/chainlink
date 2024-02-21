// Package gethwrappers_ccip provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package rebalancer

// Rebalancer contracts
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/Rebalancer/Rebalancer.abi ../../../contracts/solc/v0.8.19/Rebalancer/Rebalancer.bin Rebalancer rebalancer
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ArbitrumL1BridgeAdapter/ArbitrumL1BridgeAdapter.abi ../../../contracts/solc/v0.8.19/ArbitrumL1BridgeAdapter/ArbitrumL1BridgeAdapter.bin ArbitrumL1BridgeAdapter arbitrum_l1_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ArbitrumL2BridgeAdapter/ArbitrumL2BridgeAdapter.abi ../../../contracts/solc/v0.8.19/ArbitrumL2BridgeAdapter/ArbitrumL2BridgeAdapter.bin ArbitrumL2BridgeAdapter arbitrum_l2_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/OptimismL1BridgeAdapter/OptimismL1BridgeAdapter.abi ../../../contracts/solc/v0.8.19/OptimismL1BridgeAdapter/OptimismL1BridgeAdapter.bin OptimismL1BridgeAdapter optimism_l1_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/OptimismL2BridgeAdapter/OptimismL2BridgeAdapter.abi ../../../contracts/solc/v0.8.19/OptimismL2BridgeAdapter/OptimismL2BridgeAdapter.bin OptimismL2BridgeAdapter optimism_l2_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/NoOpOCR3/NoOpOCR3.abi ../../../contracts/solc/v0.8.19/NoOpOCR3/NoOpOCR3.bin NoOpOCR3 no_op_ocr3
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/MockBridgeAdapter/MockL2BridgeAdapter.abi ../../../contracts/solc/v0.8.19/MockBridgeAdapter/MockL2BridgeAdapter.bin MockL2BridgeAdapter mock_l2_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/MockBridgeAdapter/MockL1BridgeAdapter.abi ../../../contracts/solc/v0.8.19/MockBridgeAdapter/MockL1BridgeAdapter.bin MockL1BridgeAdapter mock_l1_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/RebalancerReportEncoder/RebalancerReportEncoder.abi ../../../contracts/solc/v0.8.19/RebalancerReportEncoder/RebalancerReportEncoder.bin RebalancerReportEncoder rebalancer_report_encoder

// Arbitrum helpers
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IArbSys/IArbSys.abi ../../../contracts/solc/v0.8.19/IArbSys/IArbSys.bin ArbSys arbsys
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/INodeInterface/INodeInterface.abi ../../../contracts/solc/v0.8.19/INodeInterface/INodeInterface.bin NodeInterface arb_node_interface
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IL2ArbitrumGateway/IL2ArbitrumGateway.abi ../../../contracts/solc/v0.8.19/IL2ArbitrumGateway/IL2ArbitrumGateway.bin L2ArbitrumGateway l2_arbitrum_gateway
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IL2ArbitrumMessenger/IL2ArbitrumMessenger.abi ../../../contracts/solc/v0.8.19/IL2ArbitrumMessenger/IL2ArbitrumMessenger.bin L2ArbitrumMessenger l2_arbitrum_messenger
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IArbRollupCore/IArbRollupCore.abi ../../../contracts/solc/v0.8.19/IArbRollupCore/IArbRollupCore.bin ArbRollupCore arbitrum_rollup_core
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IArbitrumL1GatewayRouter/IArbitrumL1GatewayRouter.abi ../../../contracts/solc/v0.8.19/IArbitrumL1GatewayRouter/IArbitrumL1GatewayRouter.bin ArbitrumL1GatewayRouter arbitrum_l1_gateway_router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IArbitrumInbox/IArbitrumInbox.abi ../../../contracts/solc/v0.8.19/IArbitrumInbox/IArbitrumInbox.bin ArbitrumInbox arbitrum_inbox
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IArbitrumGatewayRouter/IArbitrumGatewayRouter.abi ../../../contracts/solc/v0.8.19/IArbitrumGatewayRouter/IArbitrumGatewayRouter.bin ArbitrumGatewayRouter arbitrum_gateway_router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IArbitrumTokenGateway/IArbitrumTokenGateway.abi ../../../contracts/solc/v0.8.19/IArbitrumTokenGateway/IArbitrumTokenGateway.bin ArbitrumTokenGateway arbitrum_token_gateway
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/IAbstractArbitrumTokenGateway/IAbstractArbitrumTokenGateway.abi ../../../contracts/solc/v0.8.19/IAbstractArbitrumTokenGateway/IAbstractArbitrumTokenGateway.bin AbstractArbitrumTokenGateway abstract_arbitrum_token_gateway

// Generate mocks for tests
//go:generate mockery --quiet --dir ./generated/arbitrum_l1_bridge_adapter/ --name ArbitrumL1BridgeAdapterInterface --output ./mocks/mock_arbitrum_l1_bridge_adapter --outpkg mock_arbitrum_l1_bridge_adapter --case=underscore
//go:generate mockery --quiet --dir ./generated/arbitrum_l2_bridge_adapter/ --name ArbitrumL2BridgeAdapterInterface --output ./mocks/mock_arbitrum_l2_bridge_adapter --outpkg mock_arbitrum_l2_bridge_adapter --case=underscore
//go:generate mockery --quiet --dir ./generated/arbitrum_gateway_router/ --name ArbitrumGatewayRouterInterface --output ./mocks/mock_arbitrum_gateway_router --outpkg mock_arbitrum_gateway_router --case=underscore
//go:generate mockery --quiet --dir ./generated/arbitrum_inbox/ --name ArbitrumInboxInterface --output ./mocks/mock_arbitrum_inbox --outpkg mock_arbitrum_inbox --case=underscore
//go:generate mockery --quiet --dir ./generated/l2_arbitrum_gateway/ --name L2ArbitrumGatewayInterface --output ./mocks/mock_l2_arbitrum_gateway --outpkg mock_l2_arbitrum_gateway --case=underscore
//go:generate mockery --quiet --dir ./generated/arbsys/ --name ArbSysInterface --output ./mocks/mock_arbsys --outpkg mock_arbsys --case=underscore
//go:generate mockery --quiet --dir ./generated/arb_node_interface/ --name NodeInterfaceInterface --output ./mocks/mock_node_interface --outpkg mock_node_interface --case=underscore
//go:generate mockery --quiet --dir ./generated/l2_arbitrum_messenger/ --name L2ArbitrumMessengerInterface --output ./mocks/mock_l2_arbitrum_messenger --outpkg mock_l2_arbitrum_messenger --case=underscore
//go:generate mockery --quiet --dir ./generated/arbitrum_rollup_core/ --name ArbRollupCoreInterface --output ./mocks/mock_arbitrum_rollup_core --outpkg mock_arbitrum_rollup_core --case=underscore
