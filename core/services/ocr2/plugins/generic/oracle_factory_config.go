package generic

import (
	"context"
	"fmt"
	"math/big"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
)

const (
	localCfgKeyOCRContractAddress = "ocr_contract_address"
	localCfgKeyChainID             = "chain_id"
)

type ResolveOracleFactoryConfigParams struct {
	Context              context.Context
	Config               job.OracleFactoryConfig
	OnchainSigning       job.OnchainSigningStrategy
	CapabilityID         string
	LocalCfg             coreconfig.LocalCapabilities
	DefaultBootstrappers []ocrcommontypes.BootstrapperLocator
	OCRKeyBundle         ocr2key.KeyBundle
	OCRKeystore          keystore.OCR2
	EthKeystore          keystore.Eth
	Logger               logger.Logger
}

// ResolveOracleFactoryConfig fills missing oracle factory fields from node TOML config
// and local keystore defaults. Job spec values take precedence when set.
func ResolveOracleFactoryConfig(params ResolveOracleFactoryConfigParams) (job.OracleFactoryConfig, job.OnchainSigningStrategy, error) {
	cfg := params.Config
	signing := params.OnchainSigning

	if params.LocalCfg != nil && params.CapabilityID != "" {
		if capCfg := params.LocalCfg.GetCapabilityConfig(params.CapabilityID); capCfg != nil {
			local := capCfg.Config()
			if cfg.OCRContractAddress == "" {
				cfg.OCRContractAddress = local[localCfgKeyOCRContractAddress]
			}
			if cfg.ChainID == "" {
				cfg.ChainID = local[localCfgKeyChainID]
			}
		}
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
