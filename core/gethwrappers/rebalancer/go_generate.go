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
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/RebalancerReportEncoder/RebalancerReportEncoder.abi ../../../contracts/solc/v0.8.19/RebalancerReportEncoder/RebalancerReportEncoder.bin RebalancerReportEncoder rebalancer_report_encoder
