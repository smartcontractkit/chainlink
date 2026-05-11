package ocr3_1

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

// ComputeOCR3_1Config generates a full OCR3_1 configuration (signers, transmitters,
// offchain config, etc.) from oracle config parameters and node information
// fetched from Job Distributor.
func ComputeOCR3_1Config(
	env cldf.Environment,
	registryChainSel uint64,
	registryNodes []capabilities_registry_v2.INodeInfoProviderNodeInfo,
	oracleOffchainConfig V3_1OracleConfig,
	reportingPluginConfig []byte,
	extraSignerFamilies []string,
) (*ocr3.OCR2OracleConfig, error) {
	p2pIDStrings := ocr3.RegistryNodesToP2PIDs(registryNodes)

	jdNodes, err := deployment.NodeInfo(p2pIDStrings, env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info from Job Distributor: %w", err)
	}

	config, err := GenerateOCR3_1ConfigFromNodes(
		oracleOffchainConfig,
		jdNodes,
		registryChainSel,
		env.OCRSecrets,
		reportingPluginConfig,
		extraSignerFamilies,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OCR3_1 config: %w", err)
	}

	return &config, nil
}

// ComputeDKGConfig generates a full DKG OCR3_1 configuration (signers, transmitters,
// offchain config, etc.) from oracle config parameters and node information
// fetched from Job Distributor. The DKG plugin config is taken from
// oracleOffchainConfig.DKGOffchainConfig.
func ComputeDKGConfig(
	env cldf.Environment,
	registryChainSel uint64,
	registryNodes []capabilities_registry_v2.INodeInfoProviderNodeInfo,
	oracleOffchainConfig V3_1OracleConfig,
	extraSignerFamilies []string,
) (*ocr3.OCR2OracleConfig, error) {
	if oracleOffchainConfig.DKGOffchainConfig == nil {
		return nil, errors.New("DKGOffchainConfig is required for DKG configs")
	}

	p2pIDStrings := ocr3.RegistryNodesToP2PIDs(registryNodes)

	jdNodes, err := deployment.NodeInfo(p2pIDStrings, env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info from Job Distributor: %w", err)
	}

	config, err := GenerateDKGConfigFromNodes(
		oracleOffchainConfig,
		jdNodes,
		registryChainSel,
		env.OCRSecrets,
		extraSignerFamilies,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DKG config: %w", err)
	}

	return &config, nil
}
