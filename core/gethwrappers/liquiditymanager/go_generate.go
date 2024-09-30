// Package gethwrappers_ccip provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package liquiditymanager

// LiquidityManager contracts
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/LiquidityManager/LiquidityManager.abi ../../../contracts/solc/v0.8.24/LiquidityManager/LiquidityManager.bin LiquidityManager liquiditymanager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/ArbitrumL1BridgeAdapter/ArbitrumL1BridgeAdapter.abi ../../../contracts/solc/v0.8.24/ArbitrumL1BridgeAdapter/ArbitrumL1BridgeAdapter.bin ArbitrumL1BridgeAdapter arbitrum_l1_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/ArbitrumL2BridgeAdapter/ArbitrumL2BridgeAdapter.abi ../../../contracts/solc/v0.8.24/ArbitrumL2BridgeAdapter/ArbitrumL2BridgeAdapter.bin ArbitrumL2BridgeAdapter arbitrum_l2_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/OptimismL1BridgeAdapter/OptimismL1BridgeAdapter.abi ../../../contracts/solc/v0.8.24/OptimismL1BridgeAdapter/OptimismL1BridgeAdapter.bin OptimismL1BridgeAdapter optimism_l1_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/OptimismL2BridgeAdapter/OptimismL2BridgeAdapter.abi ../../../contracts/solc/v0.8.24/OptimismL2BridgeAdapter/OptimismL2BridgeAdapter.bin OptimismL2BridgeAdapter optimism_l2_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/NoOpOCR3/NoOpOCR3.abi ../../../contracts/solc/v0.8.24/NoOpOCR3/NoOpOCR3.bin NoOpOCR3 no_op_ocr3
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MockBridgeAdapter/MockL2BridgeAdapter.abi ../../../contracts/solc/v0.8.24/MockBridgeAdapter/MockL2BridgeAdapter.bin MockL2BridgeAdapter mock_l2_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MockBridgeAdapter/MockL1BridgeAdapter.abi ../../../contracts/solc/v0.8.24/MockBridgeAdapter/MockL1BridgeAdapter.bin MockL1BridgeAdapter mock_l1_bridge_adapter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/ReportEncoder/ReportEncoder.abi ../../../contracts/solc/v0.8.24/ReportEncoder/ReportEncoder.bin ReportEncoder report_encoder
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/OptimismL1BridgeAdapterEncoder/OptimismL1BridgeAdapterEncoder.abi ../../../contracts/solc/v0.8.24/OptimismL1BridgeAdapterEncoder/OptimismL1BridgeAdapterEncoder.bin OptimismL1BridgeAdapterEncoder optimism_l1_bridge_adapter_encoder

// Arbitrum helpers
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IArbSys/IArbSys.abi ../../../contracts/solc/v0.8.24/IArbSys/IArbSys.bin ArbSys arbsys
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/INodeInterface/INodeInterface.abi ../../../contracts/solc/v0.8.24/INodeInterface/INodeInterface.bin NodeInterface arb_node_interface
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IL2ArbitrumGateway/IL2ArbitrumGateway.abi ../../../contracts/solc/v0.8.24/IL2ArbitrumGateway/IL2ArbitrumGateway.bin L2ArbitrumGateway l2_arbitrum_gateway
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IL2ArbitrumMessenger/IL2ArbitrumMessenger.abi ../../../contracts/solc/v0.8.24/IL2ArbitrumMessenger/IL2ArbitrumMessenger.bin L2ArbitrumMessenger l2_arbitrum_messenger
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IArbRollupCore/IArbRollupCore.abi ../../../contracts/solc/v0.8.24/IArbRollupCore/IArbRollupCore.bin ArbRollupCore arbitrum_rollup_core
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IArbitrumL1GatewayRouter/IArbitrumL1GatewayRouter.abi ../../../contracts/solc/v0.8.24/IArbitrumL1GatewayRouter/IArbitrumL1GatewayRouter.bin ArbitrumL1GatewayRouter arbitrum_l1_gateway_router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IArbitrumInbox/IArbitrumInbox.abi ../../../contracts/solc/v0.8.24/IArbitrumInbox/IArbitrumInbox.bin ArbitrumInbox arbitrum_inbox
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IArbitrumGatewayRouter/IArbitrumGatewayRouter.abi ../../../contracts/solc/v0.8.24/IArbitrumGatewayRouter/IArbitrumGatewayRouter.bin ArbitrumGatewayRouter arbitrum_gateway_router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IArbitrumTokenGateway/IArbitrumTokenGateway.abi ../../../contracts/solc/v0.8.24/IArbitrumTokenGateway/IArbitrumTokenGateway.bin ArbitrumTokenGateway arbitrum_token_gateway
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IAbstractArbitrumTokenGateway/IAbstractArbitrumTokenGateway.abi ../../../contracts/solc/v0.8.24/IAbstractArbitrumTokenGateway/IAbstractArbitrumTokenGateway.bin AbstractArbitrumTokenGateway abstract_arbitrum_token_gateway

// Optimism helpers
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismPortal/IOptimismPortal.abi ../../../contracts/solc/v0.8.24/IOptimismPortal/IOptimismPortal.bin OptimismPortal optimism_portal
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismL2OutputOracle/IOptimismL2OutputOracle.abi ../../../contracts/solc/v0.8.24/IOptimismL2OutputOracle/IOptimismL2OutputOracle.bin OptimismL2OutputOracle optimism_l2_output_oracle
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismL2ToL1MessagePasser/IOptimismL2ToL1MessagePasser.abi ../../../contracts/solc/v0.8.24/IOptimismL2ToL1MessagePasser/IOptimismL2ToL1MessagePasser.bin OptimismL2ToL1MessagePasser optimism_l2_to_l1_message_passer
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismCrossDomainMessenger/IOptimismCrossDomainMessenger.abi ../../../contracts/solc/v0.8.24/IOptimismCrossDomainMessenger/IOptimismCrossDomainMessenger.bin OptimismCrossDomainMessenger optimism_cross_domain_messenger
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismPortal2/IOptimismPortal2.abi ../../../contracts/solc/v0.8.24/IOptimismPortal2/IOptimismPortal2.bin OptimismPortal2 optimism_portal_2
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismDisputeGameFactory/IOptimismDisputeGameFactory.abi ../../../contracts/solc/v0.8.24/IOptimismDisputeGameFactory/IOptimismDisputeGameFactory.bin OptimismDisputeGameFactory optimism_dispute_game_factory
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismStandardBridge/IOptimismStandardBridge.abi ../../../contracts/solc/v0.8.24/IOptimismStandardBridge/IOptimismStandardBridge.bin OptimismStandardBridge optimism_standard_bridge
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/IOptimismL1StandardBridge/IOptimismL1StandardBridge.abi ../../../contracts/solc/v0.8.24/IOptimismL1StandardBridge/IOptimismL1StandardBridge.bin OptimismL1StandardBridge optimism_l1_standard_bridge
