package generic

import (
	"context"
	"fmt"
	"math/big"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
)

type ResolveOracleFactoryConfigParams struct {
	Context        context.Context
	Config         job.OracleFactoryConfig
	OnchainSigning job.OnchainSigningStrategy
	// CapRegistryAddress and CapRegistryChainID are the Capabilities Registry
	// contract address and its chain ID. The oracle factory reads its OCR config
	// (signers/transmitters) from the registry, so the contract/chain default to
	// the registry's own address/chain when not set in the job spec.
	CapRegistryAddress   string
	CapRegistryChainID   string
	DefaultBootstrappers []ocrcommontypes.BootstrapperLocator
	OCRKeyBundle         ocr2key.KeyBundle
	OCRKeystore          keystore.OCR2
	EthKeystore          keystore.Eth
	Logger               logger.Logger
}

// ResolveOracleFactoryConfig fills missing oracle factory fields. Contract address
// and chain ID default to the Capabilities Registry's address/chain; the signing key
// and transmitter default to local keystore values. Job spec values take precedence
// when set.
func ResolveOracleFactoryConfig(params ResolveOracleFactoryConfigParams) (job.OracleFactoryConfig, job.OnchainSigningStrategy, error) {
	cfg := params.Config
	signing := params.OnchainSigning

	if cfg.OCRContractAddress == "" {
		cfg.OCRContractAddress = params.CapRegistryAddress
	}
	if cfg.ChainID == "" {
		cfg.ChainID = params.CapRegistryChainID
	}

	if cfg.OCRKeyBundleID == "" && params.OCRKeyBundle != nil {
		cfg.OCRKeyBundleID = params.OCRKeyBundle.ID()
	}

	if len(signing.Config) == 0 && cfg.OCRKeyBundleID != "" {
		if signing.StrategyName == "" {
			signing.StrategyName = "multi-chain"
		}
		signing.Config = map[string]string{"evm": cfg.OCRKeyBundleID}
	}

	if cfg.TransmitterID == "" && params.EthKeystore != nil && cfg.ChainID != "" {
		transmitter, err := defaultTransmitterForChain(params.Context, params.EthKeystore, cfg.ChainID)
		if err != nil {
			return cfg, signing, err
		}
		cfg.TransmitterID = transmitter
	}

	return cfg, signing, nil
}

func defaultTransmitterForChain(ctx context.Context, ethKS keystore.Eth, chainID string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	chainIDBig, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return "", fmt.Errorf("invalid chain_id %q", chainID)
	}

	addr, err := ethKS.GetRoundRobinAddress(ctx, chainIDBig)
	if err != nil {
		return "", fmt.Errorf("failed to get transmitter for chain_id %s: %w", chainID, err)
	}

	return addr.String(), nil
}

func resolveBootstrapPeers(
	specPeers []string,
	defaultBootstrappers []ocrcommontypes.BootstrapperLocator,
) ([]ocrcommontypes.BootstrapperLocator, error) {
	return ocrcommon.GetValidatedBootstrapPeers(specPeers, defaultBootstrappers, false)
}
