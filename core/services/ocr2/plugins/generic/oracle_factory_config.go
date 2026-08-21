package generic

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

type ResolveOracleFactoryConfigParams struct {
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
	// OCRContractConfig is the on-chain OCR config for this capability, read from the
	// Capabilities Registry. When set, the transmitter is taken from the entry paired
	// with this node's signer (OCRKeyBundle) rather than a round-robin keystore default,
	// so the node uses exactly the transmitter the OCR config expects. It is nil when no
	// registry config is available yet, in which case the local defaults are used.
	OCRContractConfig *ocrtypes.ContractConfig
	Logger            logger.Logger
}

// ResolveOracleFactoryConfig fills missing oracle factory fields. Contract address
// and chain ID default to the Capabilities Registry's address/chain. The signing key
// defaults to this node's OCR key bundle. The transmitter is taken from the on-chain
// OCR config entry paired with this node's signer when available, falling back to a
// round-robin keystore address otherwise. Job spec values take precedence when set.
func ResolveOracleFactoryConfig(ctx context.Context, params ResolveOracleFactoryConfigParams) (job.OracleFactoryConfig, job.OnchainSigningStrategy, error) {
	cfg := params.Config
	signing := params.OnchainSigning

	// Nothing to resolve when the oracle factory is not used by this job. Resolving would
	// otherwise force a transmitter/key lookup for jobs that never build an oracle.
	if !cfg.Enabled {
		return cfg, signing, nil
	}

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

	// Prefer the transmitter the on-chain OCR config pairs with this node's signer.
	if cfg.TransmitterID == "" && params.OCRKeyBundle != nil && params.OCRContractConfig != nil {
		if transmitter, ok := transmitterForSigner(*params.OCRContractConfig, params.OCRKeyBundle.PublicKey()); ok {
			cfg.TransmitterID = transmitter
		} else if params.Logger != nil {
			params.Logger.Warnw("node signer not found in on-chain OCR config; falling back to round-robin transmitter",
				"ocrKeyBundleID", cfg.OCRKeyBundleID)
		}
	}

	if cfg.TransmitterID == "" && params.EthKeystore != nil && cfg.ChainID != "" {
		transmitter, err := defaultTransmitterForChain(ctx, params.EthKeystore, cfg.ChainID)
		if err != nil {
			return cfg, signing, err
		}
		cfg.TransmitterID = transmitter
	}

	return cfg, signing, nil
}

// SelectOCRKeyBundleForConfig returns the key bundle whose signer public key appears in
// the on-chain OCR config. It disambiguates when a node holds multiple EVM OCR key
// bundles by picking the one the registry registered as a signer. Returns false when no
// config is provided or none of the bundles match.
func SelectOCRKeyBundleForConfig(bundles []ocr2key.KeyBundle, cc *ocrtypes.ContractConfig) (ocr2key.KeyBundle, bool) {
	if cc == nil {
		return nil, false
	}
	for _, kb := range bundles {
		pub := kb.PublicKey()
		for _, s := range cc.Signers {
			if bytes.Equal(s, pub) {
				return kb, true
			}
		}
	}
	return nil, false
}

// transmitterForSigner returns the transmitter account paired with the given signer in
// the on-chain OCR config. Signers[i] and Transmitters[i] describe the same oracle, so
// locating this node's signer yields the transmitter the OCR config expects for it.
func transmitterForSigner(cc ocrtypes.ContractConfig, signer ocrtypes.OnchainPublicKey) (string, bool) {
	for i, s := range cc.Signers {
		if bytes.Equal(s, signer) {
			if i < len(cc.Transmitters) {
				return string(cc.Transmitters[i]), true
			}
			return "", false
		}
	}
	return "", false
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
